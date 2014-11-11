package main

import "flag"
import "fmt"
import "net/http"
import "bufio"
import "os"

func getUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Sprintf("Error getting %s: %s", url, err)
	}

	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	if err != nil {
		fmt.Sprintf("Error initialize scanner for %s: %s", url, err)
	}

	r.ReadByte()
}

func main() {
	urlFilePath := flag.String("file", "urls.csv", "path to file containing urls to warm")
	flag.Parse()

	fmt.Println("Loading urls from: ", *urlFilePath)

	urlFile, err := os.Open(*urlFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer urlFile.Close()

	scanner := bufio.NewScanner(urlFile)
	for scanner.Scan() {
		currUrl := scanner.Text()
		fmt.Println("Warming cache with url: ", currUrl)

		go getUrl(currUrl)
	}
}
