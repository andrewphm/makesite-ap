package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"html/template"
)

type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content	     string
}

func main() {
	textContent, err := ioutil.ReadFile("first-post.txt")
	if err != nil {
		fmt.Println("Error reading file.")
		panic(err)
	}

	page := Page{
		TextFilePath: "first-post.txt",
		TextFileName: "first-post",
		HTMLPagePath: "first-post.html",
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
