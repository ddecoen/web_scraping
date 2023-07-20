package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocolly/colly/v2"
)

func TestCollyCollector(t *testing.T) {
	// Create a new test server to mock the responses
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a simple HTML page for testing
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
			<body>
				<h1>Test Page</h1>
				<a href="/link1">Link 1</a>
				<a href="/link2">Link 2</a>
			</body>
			</html>
		`))
	}))
	defer testServer.Close()

	// Create a new Colly collector
	c := colly.NewCollector()

	// Create a slice to store the scraped links
	var links []string

	// Set up the callback to extract links from the test page
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})

	// Start the crawling process with the test server URL
	err := c.Visit(testServer.URL)
	if err != nil {
		t.Fatalf("Crawling failed: %v", err)
	}

	// Assert that the links have been scraped correctly
	expectedLinks := []string{"/link1", "/link2"}
	for i, link := range links {
		if link != expectedLinks[i] {
			t.Errorf("Unexpected link. Expected: %s, Got: %s", expectedLinks[i], link)
		}
	}
}
