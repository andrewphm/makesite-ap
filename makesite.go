package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"html/template"
	"flag"
	"strings"
)

type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content	     string
}

func main() {
	fileNameFlag := flag.String("file", "first-post.txt", "The name of the file to convert to HTML.")

	flag.Parse()
	fileNameWithoutExtension := strings.TrimSuffix(*fileNameFlag, ".txt")

	textContent, err := ioutil.ReadFile(*fileNameFlag)
	if err != nil {
		fmt.Println("Error reading file.")
		panic(err)
	}

	page := Page{
		TextFilePath: *fileNameFlag,
		TextFileName: fileNameWithoutExtension,
		HTMLPagePath: fileNameWithoutExtension + ".html",
		Content: string(textContent),
	}

	t, err := template.ParseFiles("template.tmpl")
	if err != nil {
		fmt.Println("Error parsing template.")
		panic(err)
	}

	file, err := os.Create(page.HTMLPagePath)
	if err != nil {
		fmt.Println("Error creating file.")
		panic(err)
	}
	defer file.Close()

	err = t.Execute(file, page)
	if err != nil {
		fmt.Println("Error executing template.")
		panic(err)
	}

	// print to Stdout
	err = t.Execute(os.Stdout, page)
	if err != nil {
		fmt.Println("Error executing template.")
		panic(err)
	}

}
