package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	var threadURL string
	resp, err := http.Get(os.Args[1])
	urlArray := make([]string, 1)
	s, _ := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(s)))

	threadURL = os.Args[1]
	DirName := path.Base(threadURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	if resp.StatusCode == 200 {
		makeDir(DirName)

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "class" {
						for _, y := range n.Attr {
							if y.Key == "href" {
								if strings.Contains(y.Val, "i.4cdn.org") {
									urlArray = append(urlArray, strings.Replace(y.Val, "//", "http://", 1))
								}
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
		for i := 1; i < len(urlArray); i++ {
			fmt.Println(urlArray[i])
			downloadFile(urlArray[i], DirName)
		}
	}
}

func makeDir(s string) {
	os.Mkdir(s, os.FileMode(0777))
}

func downloadFile(picURL string, dir string) {
	fileURL, err := url.Parse(picURL)

	if err != nil {
		fmt.Println(err.Error())
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[2]
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))

	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(picURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s with %v bytes downloaded\n", fileName, size)
}
