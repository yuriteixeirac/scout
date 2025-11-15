package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/redis/go-redis/v9"
)

var rediscl = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
	Protocol: 2,
})

var ctx = context.Background()
var lemmatizer, _ = golem.New(en.New())
var stopwordsMap map[string]struct{}
var nonAlphaRegex, _ = regexp.Compile(`[^a-zA-Z]+`)

func init() {
	// Reading and/or caching of stop words
	result, _ := rediscl.SMembers(ctx, "stop_words").Result()
	if len(result) == 0 {
		content, err := os.ReadFile("data/stopwords.json")
		if err != nil {
			panic(err)
		}

		json.Unmarshal(content, &result)

		for _, word := range result {
			rediscl.SAdd(ctx, "stop_words", word)
		}
	}
	stopwordsMap = make(map[string]struct{}, len(result))
	for _, w := range result {
		stopwordsMap[w] = struct{}{}
	}
}

func ProcessWord(word string) string {
	word = nonAlphaRegex.ReplaceAllString(word, "")
	word = strings.ToLower(word)

	if _, ok := stopwordsMap[word]; ok {
		return ""
	}

	return lemmatizer.Lemma(word)
}

func TokenizeResponse(url string) error {
	wg := sync.WaitGroup{}
	// GETs the page
	req, err := http.Get(url)
	if err != nil {
		return err
	}

	defer req.Body.Close()

	if req.StatusCode > 400 {
		return fmt.Errorf("request obtained %d status code", req.StatusCode)
	}

	// Fetchs page content
	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		return err
	}

	// Remove all non-text content
	doc.Find("script, style, nav, footer, head, meta, link, noscript").Remove()

	tokenizedContent := strings.Fields(doc.Text())

	// For each token, process and caches it
	for _, word := range tokenizedContent {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			processedWord := ProcessWord(word)

			if processedWord != "" {
				rediscl.SAdd(ctx, processedWord, url)
			}
		}(word)
	}
	wg.Wait()

	return nil
}

func FetchUrls(query string) []string {
	urls := []string{}

	tokenizedQuery := strings.Fields(query)
	for _, token := range tokenizedQuery {
		currentUrls, _ := rediscl.SMembers(ctx, ProcessWord(token)).Result()
		if len(currentUrls) == 0 {
			continue
		}

		for _, url := range currentUrls {
			if !slices.Contains(urls, url) {
				urls = append(urls, url)
			}
		}
	}

	return urls
}
