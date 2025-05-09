package robotstxt

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	Client *http.Client
}

func NewParser(client *http.Client) *Parser {
	return &Parser{
		Client: client,
	}
}

type directives_t map[string]any
type agent_t string

type Robotstxt struct {
	Groups   map[agent_t]directives_t
	Sitemaps []string
}

func newRobotstxt() *Robotstxt {
	return &Robotstxt{
		Groups:   make(map[agent_t]directives_t),
		Sitemaps: []string{},
	}
}

func (self *Parser) FetchAndParse(rawURL string) (*Robotstxt, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	u.Path = "/robots.txt"

	resp, err := self.Client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("robots.txt not found or inaccessible")
	}

	scanner := bufio.NewScanner(resp.Body)
	robots := newRobotstxt()
	var current agent_t

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			current = ""
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch key {
		case "user-agent":
			current = agent_t(val)

			if _, ok := robots.Groups[current]; !ok {
				robots.Groups[current] = make(directives_t)
			}

		case "crawl-delay":
			d, _ := strconv.Atoi(val)
			robots.Groups[current]["crawl-delay"] = uint64(time.Duration(d) * time.Second)

		case "allow", "disallow", "host":
			if current == "" {
				continue
			}
			robots.Groups[current][key] = val

		case "sitemap":
			robots.Sitemaps = append(robots.Sitemaps, val)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return robots, nil
}

func (self *Robotstxt) GetRule(agent string, key string) (any, bool) {
	a := agent_t(agent)

	_, ok := self.Groups[a]
	if !ok {
		a = agent_t("*")
	}

	v, ok := self.Groups[a][key]
	if !ok {
		return nil, false
	}

	return v, true
}
