package main

import "flag"
import "fmt"
import "net/http"
import "bufio"
import "os"

func HitUrl(url string) (byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Sprintf("Error getting %s: %s", url, err)
	}

	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	if err != nil {
		fmt.Sprintf("Error initialize scanner for %s: %s", url, err)
	}

	return r.ReadByte()
}

func GetUrlBatch(r *bufio.Scanner, size int) ([]string, int) {
	var batch = make([]string, 0)
	length := 0

	for length < size {
		if r.Scan() {
			batch = append(batch, r.Text())
		} else {
			break
		}
		length++
	}

	return batch, length
}

func main() {
	urlFilePath := flag.String("file", "urls.csv", "path to file containing urls to warm")
	flag.Parse()

	fmt.Println("Loading urls from: ", *urlFilePath)

	urlFile, err := os.Open(*urlFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}

	scanner := bufio.NewScanner(urlFile)

	currBatch, length := GetUrlBatch(scanner, 2)
	for i := 0; i < length; i++ {
		fmt.Println("Warming cache with url: ", currBatch[i])
		go HitUrl(currBatch[i])
	}
	defer urlFile.Close()
}
