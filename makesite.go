package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
}

func main() {
	dirFlag := flag.String("dir", "", "The directory to search for text files")
	fileFlag := flag.String("file", "", "The single text file to process")
	flag.Parse()


	htmlDir := "html_pages"
	if err := os.MkdirAll(htmlDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", htmlDir, err)
		panic(err)
	}

	if *fileFlag != "" {
		processFile(*fileFlag, htmlDir)
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
				processFile(filePath, htmlDir)
			}
		}
	} else {
		fmt.Println("No file or directory specified. Use --file or --dir flag.")
	}
}

func processFile(filePath, htmlDir string) {
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ".txt")
	htmlFilePath := filepath.Join(htmlDir, fileNameWithoutExt+".html")

	page := Page{
		TextFilePath: filePath,
		TextFileName: fileNameWithoutExt,
		HTMLPagePath: htmlFilePath,
		Content:      readFileContent(filePath),
	}

	generateHTMLPage(page)
}

func readFileContent(filePath string) string {
	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return ""
	}
	return string(contentBytes)
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
