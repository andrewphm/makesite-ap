package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/russross/blackfriday/v2"
)

type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      template.HTML
}

func main() {
	startTime := time.Now()

	dirFlag := flag.String("dir", "", "The directory to search for text files")
	fileFlag := flag.String("file", "", "The single text file to process")
	flag.Parse()

	pageCount := 0
	totalSizeBytes := int64(0)

	htmlDir := "html_pages"
	if err := os.MkdirAll(htmlDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", htmlDir, err)
		panic(err)
	}

	if *fileFlag != "" {
		processFile(*fileFlag, htmlDir, &pageCount, &totalSizeBytes)
	} else if *dirFlag != "" {
		files, err := ioutil.ReadDir(*dirFlag)
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", *dirFlag, err)
			panic(err)
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
				fmt.Println(file.Name())
				filePath := filepath.Join(*dirFlag, file.Name())
				processFile(filePath, htmlDir, &pageCount, &totalSizeBytes)
			}
		}
	} else {
		fmt.Println("No file or directory specified. Use --file or --dir flag.")
	}

	elapsed := time.Since(startTime)
	totalSizeKB := float64(totalSizeBytes) / 1024.0

	fmt.Printf("\033[1m\033[32mSuccess!\033[0m Generated \033[1m%d\033[0m pages.\n", pageCount)
	fmt.Printf("Total size %.2fKB. Total time %.2f ms. \n", totalSizeKB, float64(elapsed.Milliseconds()))
}

func processFile(filePath, htmlDir string, pageCount *int, totalSizeBytes *int64) {
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ".txt")
	htmlFilePath := filepath.Join(htmlDir, fileNameWithoutExt+".html")

	page := Page{
		TextFilePath: filePath,
		TextFileName: fileNameWithoutExt,
		HTMLPagePath: htmlFilePath,
		Content:      template.HTML(readFileContent(filePath)),
	}

	generateHTMLPage(page)
	*pageCount++
	info, err := os.Stat(page.HTMLPagePath)
	if err == nil {
		*totalSizeBytes += info.Size()
	}
}

func readFileContent(filePath string) string {
	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return ""
	}

	output := blackfriday.Run(contentBytes)

	return string(output)
}

func generateHTMLPage(page Page) {
	t, err := template.ParseFiles("template.tmpl")
	if err != nil {
		fmt.Printf("Error parsing template for file %s: %v\n", page.TextFilePath, err)
		return
	}

	file, err := os.Create(page.HTMLPagePath)
	if err != nil {
		fmt.Printf("Error creating HTML file %s: %v\n", page.HTMLPagePath, err)
		return
	}
	defer file.Close()

	if err := t.Execute(file, page); err != nil {
		fmt.Printf("Error executing template for file %s: %v\n", page.TextFilePath, err)
	}
}
