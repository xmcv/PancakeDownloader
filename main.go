package main

import (
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	
	"golang.org/x/net/html"
)

func main() {

	//Arguments check
	if len(os.Args) == 1 {
		fmt.Printf("Usage: ./main <URL>")
		os.Exit(1)
	}
	
	resp, err := http.Get(os.Args[1])
	
	if err != nil {
		log.Fatal()
	}
	
	var urlArray []string
	s, _ := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(s)))
	dirName := path.Base(resp.Request.URL.Path)
	
	if err != nil {
		log.Fatal()
	}
	
	if resp.StatusCode == 200 {

		// Creating the directory
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			os.Mkdir(dirName, os.FileMode(0775))
		}

		// Parsing the HTML Body, looking for pics links
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "class" && a.Val == "fileThumb" {
						for _, y := range n.Attr {
							if y.Key == "href" {

								// Add http:// to the link
								urlArray = append(urlArray, strings.Replace(y.Val, "//", "http://", 1))
							}
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		// Download all files
		for i := 0; i < len(urlArray); i++ {
			fmt.Printf("[*]Downloading %d of %d: %s\n", i, len(urlArray), urlArray[i])
			downloadFile(urlArray[i], dirName)
		}

	}
}

func downloadFile(picUrl string, dir string) {
	fileName := path.Base(picUrl)

	// Checking duals
	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, fileName)); err == nil {
		fmt.Printf("File already downloaded, skipping...\n")
		return
	}

	// Downloading files in the thread's folder [dir/fileName]
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))

	if err != nil {
		log.Fatal()
	}
	defer file.Close()

	resp, err := http.Get(picUrl)

	if err != nil {
		log.Fatal()
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		log.Fatal()
	}

	fmt.Printf("%s with %v bytes downloaded\n", fileName, size)
}
