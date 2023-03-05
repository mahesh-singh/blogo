package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title   string
	Heading string
	Content string
}

type Home struct {
	Title string
	Body  string
}

type Data struct {
	Content template.HTML
	Title   string
}

func main() {
	createHomePage()
	createPostPages()

}

func createPostPages() {
	individualPostTemplate := []string{
		"./template/base.tmpl",
		"./template/posts/post.tmpl",
	}

	inFolder := "./content/posts/"
	outFolder := "./public/posts"
	files, _ := os.ReadDir(inFolder)

	//Generating individual post
	for _, file := range files {
		//if filepath.Ext(file.Name()) == ".md" {}
		markdownFile, _ := os.Open(inFolder + "/" + file.Name())

		// don't forget to close it when done
		defer markdownFile.Close()

		// create the html file
		htmlFilePath := outFolder + "/" + strings.Replace(file.Name(), ".md", "", -1) + ".html"
		if _, err := os.Stat(filepath.Dir(htmlFilePath)); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(htmlFilePath), 0770)
			if err != nil {
				log.Printf("error while creating post directory: %v", err)
			}
		}
		htmlFile, _ := os.Create(htmlFilePath)
		defer htmlFile.Close()

		// read the md
		reader := bufio.NewReader(markdownFile)
		markdown, _ := io.ReadAll(reader)

		// send the md to blackfriday
		html := blackfriday.MarkdownCommon(markdown)

		d := Data{Content: template.HTML(html), Title: strings.Replace(file.Name(), "-", " ", -1)}
		templ := template.Must(template.ParseFiles(individualPostTemplate...))
		err := templ.ExecuteTemplate(htmlFile, "base", d)
		if err != nil {
			log.Printf("error: %v", err)
		}

	}

	//Generating post index
	indexPostTemplate := []string{
		"./template/base.tmpl",
		"./template/posts/index.tmpl",
	}

	htmlFile, _ := os.Create("./public/posts/index.html")
	files, _ = os.ReadDir("./public/posts")

	htmlListContent := "<ul>"
	for _, file := range files {
		if file.Name() == "index.html" {
			continue
		}
		htmlListContent = htmlListContent + fmt.Sprintf("<li><a href='%s'/>%s</li>", file.Name(), strings.Replace(file.Name(), "-", " ", -1))
	}
	htmlListContent = htmlListContent + "</ul>"

	d := Data{Content: template.HTML(htmlListContent), Title: "Post list"}
	templ := template.Must(template.ParseFiles(indexPostTemplate...))
	err := templ.ExecuteTemplate(htmlFile, "base", d)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func createHomePage() {
	homeTemplate := []string{
		"./template/base.tmpl",
		"./template/index.tmpl",
	}
	//Rander home page
	homeContentPath := "./content/index.md"

	//homePublicFolder := "./public"
	homeMarkdownFile, _ := os.Open(homeContentPath)

	defer homeMarkdownFile.Close()

	reader := bufio.NewReader(homeMarkdownFile)
	homeMarkdown, _ := io.ReadAll(reader)

	content := blackfriday.MarkdownCommon(homeMarkdown)
	//log.Print(string(content))

	d := Data{Content: template.HTML(content), Title: "Mahesh Home Page"}

	templ := template.Must(template.ParseFiles(homeTemplate...))
	indexFile, _ := os.Create("./public/index.html")
	defer indexFile.Close()
	err := templ.ExecuteTemplate(indexFile, "base", d)
	if err != nil {
		log.Printf("error: %v", err)
	}
}
