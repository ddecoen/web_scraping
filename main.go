// Code inspired by Python_Scrapy WebCrawl Example

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Set up struct for key items from webcrawl
type WebFocusedCrawlItem struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Text  []string `json:"text"`
	Tags  []string `json:"tags"`
}

// Create function to remove stopwords
func removeStopwords(tokens []string) []string {
	stopwordList := []string{"a", "an", "the", "and", "or", "of", "to", "in", "is", "that"}
	var goodTokens []string
	for _, token := range tokens {
		found := false
		for _, stopword := range stopwordList {
			if strings.ToLower(token) == stopword {
				found = true
				break
			}
		}
		if !found {
			goodTokens = append(goodTokens, token)
		}
	}
	return goodTokens
}

func main() {
	// Time the code
	startTime := time.Now()

	// Create a new collector using colly
	c := colly.NewCollector()

	// list of Wikipedia URLs for topic of interest
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	// Create a slice to store scraped items
	var items []WebFocusedCrawlItem

	// On HTML response, perform scraping
	c.OnHTML("html", func(e *colly.HTMLElement) {
		item := WebFocusedCrawlItem{}
		item.URL = e.Request.URL.String()
		item.Title = e.ChildText("h1")
		item.Text = e.ChildTexts("div#mw-content-text p")

		// Tags
		tags := strings.Split(e.Request.URL.Path, "/")
		var tagsList []string
		for _, tag := range tags {
			tag = strings.ToLower(tag)
			tag = strings.Trim(tag, "_")
			tagsList = append(tagsList, removeStopwords(strings.Fields(tag))...)
		}
		item.Tags = tagsList

		items = append(items, item)
	})

	// On request error
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping for each URL
	for _, url := range urls {
		err := c.Visit(url)
		if err != nil {
			log.Println("Error visiting URL:", url, "\nError:", err)
		}
	}

	// Convert items to JSON
	jsonData, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		log.Fatal("Failed to convert to JSON:", err)
	}

	// Save JSON data to items.jl
	err = ioutil.WriteFile("items.jl", jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to save JSON data to items.jl:", err)
	}

	fmt.Println("Data saved to items.jl")

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	fmt.Printf("Execution time: %s\n", executionTime)
}
