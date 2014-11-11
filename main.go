package main

import "flag"
import "fmt"
import "net/http"
import "bufio"
import "os"
import "sync"
import "time"

func HitUrl(url string) (byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Sprintf("Error getting %s: %s", url, err)
	}

	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	if err != nil {
		fmt.Sprintln("Error initialize scanner for %s: %s", url, err)
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

func SetupScanner(filepath string) (*bufio.Scanner, *os.File) {
	urlFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}

	scanner := bufio.NewScanner(urlFile)

	return scanner, urlFile
}

func HitUrlsInBatches(s *bufio.Scanner, batchSize int, waitTime time.Duration) int {
	var wg sync.WaitGroup
	totalHit := 0

	currBatch, length := GetUrlBatch(s, batchSize)
	for length > 0 {
		for i := 0; i < length; i++ {
			wg.Add(1)

			fmt.Println("Warming cache with url: ", currBatch[i])
			go func(wg *sync.WaitGroup, url string) {
				HitUrl(url)
				totalHit++
				wg.Done()
			}(&wg, currBatch[i])
		}

		wg.Wait()
		currBatch, length = GetUrlBatch(s, batchSize)

		if length > 0 {
			time.Sleep(waitTime)
		}
	}

	return totalHit
}

func main() {
	urlFilePath := flag.String("file", "urls.csv", "path to file containing urls to warm")
	waitTime := flag.Duration("waitTime", 500*time.Millisecond, "time between batches of requests")
	batchSize := flag.Int("batchSize", 2, "number of urls to request in parallel")
	flag.Parse()

	fmt.Println("Loading urls from: ", *urlFilePath)
	scanner, urlFile := SetupScanner(*urlFilePath)

	totalHit := HitUrlsInBatches(scanner, *batchSize, *waitTime)

	fmt.Sprintln("Warmed cache with %i urls", totalHit)

	urlFile.Close()
}
