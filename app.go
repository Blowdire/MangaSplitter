package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// func extractImagesFromPDF(pdfPath string, outputFolder string) {
// 	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
// 		fmt.Println("Error creating output directory:", err)
// 		return
// 	}
// 	os.RemoveAll(outputFolder)
// 	os.Mkdir(outputFolder, os.ModePerm)

// 	pdfDocument, err := fitz.New(pdfPath)
// 	if err != nil {
// 		fmt.Println("Error opening PDF:", err)
// 		return
// 	}
// 	defer pdfDocument.Close()

// 	x := 1

// 	for pageNumber := 0; pageNumber < pdfDocument.NumPage(); pageNumber++ {
// 		page, err := pdfDocument.Page(pageNumber)
// 		if err != nil {
// 			fmt.Println("Error accessing page:", err)
// 			continue
// 		}
// 		images := page.GetImages()
// 		for imgIndex, img := range images {
// 			xref := img.Xref
// 			baseImage, err := pdfDocument.Image(xref)
// 			if err != nil {
// 				fmt.Println("Error extracting image:", err)
// 				continue
// 			}
// 			imageBytes := baseImage.Image

// 			reader := bytes.NewReader(imageBytes)
// 			img, _, err := image.Decode(reader)
// 			if err != nil {
// 				fmt.Println("Error decoding image:", err)
// 				continue
// 			}

// 			width := img.Bounds().Max.X
// 			height := img.Bounds().Max.Y

// 			if width > height {
// 				singlePageWidth := width / 2

// 				leftPage := resize.Resize(uint(singlePageWidth), 0, img, resize.Lanczos3)
// 				leftPagePath := filepath.Join(outputFolder, fmt.Sprintf("page_%03d_%02d.jpg", x, imgIndex))
// 				saveImage(leftPage, leftPagePath)

// 				x++
// 				pageNumber := fmt.Sprintf("%03d", x)
// 				x++

// 				rightPage := resize.Resize(uint(singlePageWidth), 0, img, resize.Lanczos3)
// 				rightPagePath := filepath.Join(outputFolder, fmt.Sprintf("page_%03d_%02d.jpg", pageNumber, imgIndex))
// 				saveImage(rightPage, rightPagePath)
// 			} else {
// 				pageNumber := fmt.Sprintf("%03d", x)
// 				pagePath := filepath.Join(outputFolder, fmt.Sprintf("page_%03d_%02d.jpg", pageNumber, imgIndex))
// 				saveImage(img, pagePath)
// 				x++
// 			}
// 		}
// 	}
// }

// func saveImage(img image.Image, outputPath string) {
// 	outputFile, err := os.Create(outputPath)
// 	if err != nil {
// 		fmt.Println("Error creating output image:", err)
// 		return
// 	}
// 	defer outputFile.Close()

// 	jpeg.Encode(outputFile, img, nil)
// }

// func compressToZip(folderPath, outputZipPath string) {
// 	outputFile, err := os.Create(outputZipPath)
// 	if err != nil {
// 		fmt.Println("Error creating output ZIP file:", err)
// 		return
// 	}
// 	defer outputFile.Close()

// 	zipWriter := zip.NewWriter(outputFile)
// 	defer zipWriter.Close()

// 	filepath.Walk(folderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
// 		if err != nil {
// 			fmt.Println("Error accessing file:", err)
// 			return nil
// 		}

// 		if !fileInfo.IsDir() {
// 			file, err := os.Open(filePath)
// 			if err != nil {
// 				fmt.Println("Error opening file:", err)
// 				return nil
// 			}
// 			defer file.Close()

// 			fileHeader, err := zip.FileInfoHeader(fileInfo)
// 			if err != nil {
// 				fmt.Println("Error creating file header:", err)
// 				return nil
// 			}

// 			fileHeader.Name = filepath.Base(filePath)
// 			writer, err := zipWriter.CreateHeader(fileHeader)
// 			if err != nil {
// 				fmt.Println("Error creating ZIP entry:", err)
// 				return nil
// 			}

// 			_, err = io.Copy(writer, file)
// 			if err != nil {
// 				fmt.Println("Error copying file to ZIP:", err)
// 			}
// 		}

// 		return nil
// 	})
// }

func (a *App) ChooseFile(updateProgress func(int)) {
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
	for i := 0; i <= 100; i += 10 {
		updateProgress(i)
		time.Sleep(1 * time.Second)
	}
}
