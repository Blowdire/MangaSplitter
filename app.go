package main

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func extractImagesFromPDF(pdfPath string, outputFolder string, a *App) {
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

	x := 1

	for pageNumber := 0; pageNumber < pdfDocument.NumPage(); pageNumber++ {
		//page, err := pdfDocument[pageNumber]
		//if err != nil {
		//	fmt.Println("Error accessing page:", err)
		//	continue
		//}
		//images := page.GetImages()
		images := make([]image.Image, 1)
		currentImage, err := pdfDocument.Image(pageNumber)
		if err != nil {
			fmt.Println("error getting image")
			continue
		}
		images[0] = currentImage

		for imgIndex, img := range images {

			width := img.Bounds().Max.X
			height := img.Bounds().Max.Y

			if width > height {
				singlePageWidth := width / 2

				leftPage := resize.Resize(uint(singlePageWidth), 0, img, resize.Lanczos3)
				leftPagePath := filepath.Join(outputFolder, fmt.Sprintf("page_%03d_%02d.jpg", x, imgIndex))
				saveImage(leftPage, leftPagePath)

				x++
				pageNumber := fmt.Sprintf("%03d", x)
				x++

				rightPage := resize.Resize(uint(singlePageWidth), 0, img, resize.Lanczos3)
				rightPagePath := filepath.Join(outputFolder, fmt.Sprintf("page_%03d_%02d.jpg", pageNumber, imgIndex))
				saveImage(rightPage, rightPagePath)
			} else {
				pageNumber := fmt.Sprintf("%03d", x)
				pagePath := filepath.Join(outputFolder, fmt.Sprintf("page_%03d_%02d.jpg", pageNumber, imgIndex))
				saveImage(img, pagePath)
				x++
			}
		}
		result := PageResult{PageNumber: pageNumber, CurrentTotalPages: pdfDocument.NumPage()}
		runtime.EventsEmit(a.ctx, "pageDone", result)
	}
}

func saveImage(img image.Image, outputPath string) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating output image:", err)
		return
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, img, nil)
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

func (a *App) ChooseFile() {
	filepath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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
	fmt.Println(filepath)
	filename := path.Base(filepath)
	extension := path.Ext(filename)
	filenameWithoutExtension := strings.TrimSuffix(filename, extension)
	fmt.Print(filenameWithoutExtension)
	// extractImagesFromPDF(filepath, "./"+filenameWithoutExtension, a)
	// compressToZip("./"+filenameWithoutExtension, "./"+filenameWithoutExtension+".cbz")
}

type PageResult struct {
	PageNumber        int
	CurrentTotalPages int
}
