package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
	lrucache "github.com/guilherme13c/go-search/utils/lru-cache"
	"github.com/guilherme13c/go-search/utils/queue"
	"github.com/guilherme13c/go-search/utils/set"
)

var (
	robotsCache = lrucache.LRUCache[string, map[string]string]{}
	robotsMu    = sync.Mutex{}
)

const (
	userAgent    = "go-search-bot/0.0.1"
	maxCacheSize = 256
)

func main() {
	os.RemoveAll("corpus/")
	os.Mkdir("corpus", 0777)

	frontier := queue.NewQueue[string]()
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

	for {
		semaphore <- struct{}{}
		go func() {
			defer func() { <-semaphore }()
			_, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			url, ok := frontier.Get()
			if !ok {
				return
			}

			domain, err := getDomain(url)
			if err != nil {
				return
			}

			{
				robotsMu.Lock()
				defer robotsMu.Unlock()

				_, exists := robotsCache.Get(domain)
				if !exists {
					respRobots, err := soup.Get(domain + "/robots.txt")
					if err != nil {
						return
					}

					print(respRobots)
				}
			}

			resp, err := http.Get(url)
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
				if visited.Contains(url) {
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
				fmt.Sprintf("WARC-Target-URI: %s", url),
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

func parseCrawlDelay(robotsTxt, userAgent string) time.Duration {
	lines := strings.Split(robotsTxt, "\n")
	matchedAgent := false
	delay := time.Duration(0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			agent := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			matchedAgent = (agent == "*" || strings.EqualFold(agent, userAgent))
		} else if matchedAgent && strings.HasPrefix(strings.ToLower(line), "crawl-delay:") {
			val := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			if secs, err := time.ParseDuration(val + "s"); err == nil {
				delay = secs
			}
			matchedAgent = false
		}
	}
	return delay
}
