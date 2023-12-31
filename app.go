package main

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/go-fitz"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func extractImagesFromPDF(pdfPath string, outputFolder string, a *App, compressionLevel int) {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}
	os.RemoveAll(outputFolder)
	os.Mkdir(outputFolder, os.ModePerm)

	pdfDocument, err := fitz.New(pdfPath)
	if err != nil {
		fmt.Println("Error opening PDF:", err)
		return
	}
	defer pdfDocument.Close()
	doneChannel := make(chan struct{}, pdfDocument.NumPage())
	numberOfSplits := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	numPages := pdfDocument.NumPage()
	for pageNumber := 0; pageNumber < numPages; pageNumber++ {
		wg.Add(1)
		go func(pageNumber int) {
			defer func() {
				doneChannel <- struct{}{}
			}()
			images := make([]image.Image, 1)
			mu.Lock()
			currentImage, err := pdfDocument.Image(pageNumber)
			mu.Unlock()
			if err != nil {
				fmt.Println("error getting image")
				return
			}
			images[0] = currentImage

			for _, img := range images {

				width := img.Bounds().Max.X
				height := img.Bounds().Max.Y

				if width > height {
					singlePageWidth := width / 2
					mu.Lock() // Lock the mutex before modifying x
					imageNumber := pageNumber + (2 + numberOfSplits)
					numberOfSplits++
					mu.Unlock()
					leftPage := imaging.Crop(img, image.Rect(0, 0, singlePageWidth, img.Bounds().Dy()))
					leftPagePath := filepath.Join(outputFolder, fmt.Sprintf("%03d_l.jpg", imageNumber))
					saveImage(leftPage, leftPagePath, compressionLevel)

					rightPage := imaging.Crop(img, image.Rect(singlePageWidth, 0, width, img.Bounds().Dy()))
					rightPagePath := filepath.Join(outputFolder, fmt.Sprintf("%03d_r.jpg", imageNumber))
					saveImage(rightPage, rightPagePath, compressionLevel)
					mu.Lock() // Lock the mutex before modifying x

				} else {
					mu.Lock() // Lock the mutex before modifying x
					imageNumber := pageNumber + (2 + numberOfSplits)

					mu.Unlock()
					pagePath := filepath.Join(outputFolder, fmt.Sprintf("%03d.jpg", imageNumber))
					saveImage(img, pagePath, compressionLevel)

				}
			}
			wg.Done()
		}(pageNumber)

	}
	go func() {
		wg.Wait()
		close(doneChannel)
	}()
	numTimeouts := 0
	for i := 0; i < pdfDocument.NumPage(); i++ {
		select {
		case <-doneChannel:
			result := PageResult{PageNumber: i + 1, CurrentTotalPages: pdfDocument.NumPage()}
			fmt.Println("Page Done", result)
			runtime.EventsEmit(a.ctx, "pageDone", result)
		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			numTimeouts++
			if numTimeouts >= 2 {
				close(doneChannel)
				break

			}
			continue
		}
	}

	fmt.Println("Done")
	result := PageResult{PageNumber: pdfDocument.NumPage(), CurrentTotalPages: pdfDocument.NumPage()}
	fmt.Println("Page Done", result)
	runtime.EventsEmit(a.ctx, "pageDone", result)

}
func saveImage(img image.Image, outputPath string, compressionLevel int) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating output image:", err)
		return
	}
	defer outputFile.Close()
	options := jpeg.Options{Quality: compressionLevel}
	jpeg.Encode(outputFile, img, &options)
}

func compressToZip(folderPath, outputZipPath string) {
	outputFile, err := os.Create(outputZipPath)
	if err != nil {
		fmt.Println("Error creating output ZIP file:", err)
		return
	}
	defer outputFile.Close()

	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	filepath.Walk(folderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing file:", err)
			return nil
		}

		if !fileInfo.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return nil
			}
			defer file.Close()

			fileHeader, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				fmt.Println("Error creating file header:", err)
				return nil
			}

			fileHeader.Name = filepath.Base(filePath)
			fileHeader.Method = zip.Deflate
			writer, err := zipWriter.CreateHeader(fileHeader)
			if err != nil {
				fmt.Println("Error creating ZIP entry:", err)
				return nil
			}

			_, err = io.Copy(writer, file)
			if err != nil {
				fmt.Println("Error copying file to ZIP:", err)
			}
		}

		return nil
	})
}

func (a *App) ChooseFile(compressionLevel int) {
	filepathSel, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a PDF file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files",
				Pattern:     "*.pdf",
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	filename := filepath.Base(filepathSel)
	extension := path.Ext(filename)
	fileDir := path.Dir(filepathSel)
	filenameWithoutExtension := strings.TrimSuffix(filename, extension)
	fmt.Print(filenameWithoutExtension)
	outTempDir := fileDir + "/" + filenameWithoutExtension
	extractImagesFromPDF(filepathSel, outTempDir, a, compressionLevel)
	compressToZip(outTempDir, outTempDir+".cbz")
	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", fileDir)
	case "windows":
		cmd = exec.Command("explorer", fileDir)
	case "darwin":
		cmd = exec.Command("open", fileDir)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return

	}
}

type PageResult struct {
	PageNumber        int
	CurrentTotalPages int
}
