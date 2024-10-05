package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	inputFile := flag.String("i", "", "Path to input HTML file")
	outputFile := flag.String("o", "", "Path to output HTML file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		log.Fatal("Both input and output file paths are required")
	}

	// Read the input HTML file
	content, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	// Process the HTML tree
	processNode(doc, filepath.Dir(*inputFile))

	// Create the output file
	outFile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()

	// Write the modified HTML to the output file
	err = html.Render(outFile, doc)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Conversion completed successfully")
}

func processNode(n *html.Node, baseDir string) {
	if n.Type == html.ElementNode && n.Data == "img" {
		for i, attr := range n.Attr {
			if attr.Key == "src" {
				if strings.HasPrefix(attr.Val, "http://") || strings.HasPrefix(attr.Val, "https://") {
					continue
				}
				imagePath := filepath.Join(baseDir, attr.Val)
				base64Data, err := imageToBase64(imagePath)
				if err != nil {
					log.Printf("Error converting image to base64: %v", err)
					continue
				}
				n.Attr[i].Val = fmt.Sprintf("data:image/png;base64,%s", base64Data)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processNode(c, baseDir)
	}
}

func imageToBase64(imagePath string) (string, error) {
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}
