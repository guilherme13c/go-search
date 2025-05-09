package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/google/uuid"

	lrucache "github.com/guilherme13c/go-search/utils/lru-cache"
	"github.com/guilherme13c/go-search/utils/queue"
	robotstxt "github.com/guilherme13c/go-search/utils/robots-txt"
	"github.com/guilherme13c/go-search/utils/set"
)

const (
	userAgent    = "go-search-bot/0.0.1"
	maxCacheSize = 1024
)

var (
	robotsCache = lrucache.NewLruCache[string, robotstxt.Robotstxt](maxCacheSize)
	robotsMu    = sync.Mutex{}
)

func main() {
	os.RemoveAll("corpus/")
	os.Mkdir("corpus", 0777)

	frontier := queue.NewQueue[string](0)
	visited := set.NewSet[string]()
	semaphore := make(chan struct{}, 512)

	seedFile, errOpenSeedFile := os.Open("crawler/seeds.txt")
	if errOpenSeedFile != nil {
		panic(errOpenSeedFile)
	}
	defer seedFile.Close()

	scanner := bufio.NewScanner(seedFile)
	for scanner.Scan() {
		frontier.Put(scanner.Text())
	}

	robotsParser := robotstxt.NewParser(&http.Client{Timeout: time.Second * 5})

	run := true

	for run {
		semaphore <- struct{}{}
		go func() {
			defer func() { <-semaphore }()
			_, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			pageUrl, ok := frontier.Get()
			if !ok {
				return
			}

			domain, err := getDomain(pageUrl)
			if err != nil {
				return
			}

			var robotsTxt robotstxt.Robotstxt
			{
				robotsMu.Lock()
				defer robotsMu.Unlock()

				robotRules, exists := robotsCache.Get(domain)
				if !exists {
					rules, err := robotsParser.FetchAndParse(domain)
					if err != nil || rules == nil {
						return
					}

					robotsCache.Put(domain, *rules)
					robotRules = rules
				}
				robotsTxt = *robotRules
			}
			fmt.Printf("%#v\n", robotsTxt)

			u, err := url.Parse(pageUrl)
			if err != nil {
				return
			}
			if u.Scheme == "" {
				u.Scheme = "https"
			}

			client := http.Client{Timeout: time.Second * 5}
			resp, err := client.Get(u.String())
			if err != nil {
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return
			}
			contentType := resp.Header.Get("Content-Type")
			if !strings.Contains(contentType, "text/html") &&
				!strings.Contains(contentType, "application/xhtml+xml") &&
				!strings.Contains(contentType, "application/xml") {
				return
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			body := string(bodyBytes)

			parsed := soup.HTMLParse(body)
			links := parsed.FindAll("a")
			for _, link := range links {
				extractedUrl, ok := link.Attrs()["href"]
				if !ok {
					continue
				}
				if strings.HasPrefix(extractedUrl, "/") || strings.HasPrefix(extractedUrl, "#") {
					extractedUrl = domain + extractedUrl
				}
				if visited.Contains(pageUrl) {
					continue
				}
				frontier.Put(extractedUrl)
			}

			docId := uuid.NewString()
			file, err := os.Create("corpus/" + docId + ".warc")
			if err != nil {
				return
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			defer writer.Flush()

			warcDate := time.Now().UTC().Format("2006-01-02T15:04:05Z")
			recordId := uuid.New().String()
			headers := []string{
				"WARC/1.0",
				"WARC-Type: response",
				fmt.Sprintf("WARC-Date: %s", warcDate),
				fmt.Sprintf("WARC-Record-ID: <urn:uuid:%s>", recordId),
				fmt.Sprintf("WARC-Target-URI: %s", pageUrl),
				"Content-Type: application/http; msgtype=response",
				fmt.Sprintf("Content-Length: %d", len(bodyBytes)+len(contentType)+50),
				"",
				fmt.Sprintf("HTTP/1.1 %s", resp.Status),
			}
			for k, v := range resp.Header {
				headers = append(headers, fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
			}
			headers = append(headers, "", body)
			for _, line := range headers {
				writer.WriteString(line + "\r\n")
			}
		}()
	}
}

func getDomain(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid URL")
	}
	return strings.Join(parts[:3], "/"), nil
}
