package robotstxt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchAndParse_Success(t *testing.T) {
	const robotsContent = `
User-agent: Googlebot
Disallow: /private
Crawl-delay: 5
Sitemap: https://example.com/sitemap.xml

User-agent: *
Allow: /
`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/robots.txt" {
			t.Errorf("expected request to /robots.txt, got %s", r.URL.Path)
		}
		fmt.Fprint(w, robotsContent)
	}))
	defer server.Close()

	parser := &Parser{Client: server.Client()}
	rtxt, err := parser.FetchAndParse(server.URL)
	if err != nil {
		t.Fatalf("FetchAndParse error = %v; want nil", err)
	}

	if got := len(rtxt.Groups); got != 2 {
		t.Errorf("len(Groups) = %d; want 2", got)
	}

	v, ok := rtxt.GetRule("Googlebot", "disallow")
	if !ok {
		t.Fatal("expected disallow rule for Googlebot, got none")
	}
	if s, _ := v.(string); s != "/private" {
		t.Errorf("Googlebot disallow = %q; want %q", s, "/private")
	}

	v, ok = rtxt.GetRule("Googlebot", "crawl-delay")
	if !ok {
		t.Fatal("expected crawl-delay for Googlebot, got none")
	}
	if d, _ := v.(uint64); d != uint64(5*time.Second) {
		t.Errorf("Googlebot crawl-delay = %v; want %v", time.Duration(d), 5*time.Second)
	}

	v, ok = rtxt.GetRule("*", "allow")
	if !ok {
		t.Fatal("expected allow rule for *, got none")
	}
	if s, _ := v.(string); s != "/" {
		t.Errorf("Wildcard allow = %q; want %q", s, "/")
	}

	if len(rtxt.Sitemaps) != 1 || rtxt.Sitemaps[0] != "https://example.com/sitemap.xml" {
		t.Errorf("Sitemaps = %v; want [%q]", rtxt.Sitemaps, "https://example.com/sitemap.xml")
	}
}

func TestFetchAndParse_InvalidURL(t *testing.T) {
	parser := &Parser{Client: http.DefaultClient}
	if _, err := parser.FetchAndParse("://%%invalid%%%"); err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestFetchAndParse_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	parser := &Parser{Client: server.Client()}
	if _, err := parser.FetchAndParse(server.URL); err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestGetRule_Fallback(t *testing.T) {
	rtxt := &Robotstxt{
		Groups: map[agent_t]directives_t{
			"*": {"disallow": "/default"},
		},
	}
	v, ok := rtxt.GetRule("NonExistentBot", "disallow")
	if !ok {
		t.Fatal("expected fallback rule for *, got none")
	}
	if s, _ := v.(string); s != "/default" {
		t.Errorf("fallback disallow = %q; want %q", s, "/default")
	}
}

func TestGetRule_MissingKey(t *testing.T) {
	rtxt := &Robotstxt{
		Groups: map[agent_t]directives_t{
			"Googlebot": {"allow": "/public"},
		},
	}
	if _, ok := rtxt.GetRule("Googlebot", "disallow"); ok {
		t.Error("expected no disallow rule for Googlebot, but got one")
	}
}
