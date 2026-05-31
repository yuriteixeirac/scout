package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/redis/go-redis/v9"
)

// ok to ignore error
var DB, _ = strconv.ParseInt(os.Getenv("REDIS_DB"), 0, 32)
var protocol, _ = strconv.ParseInt(os.Getenv("REDIS_PROTOCOL"), 0, 32)

var rediscl = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_SERVER"),
	DB:       int(DB),
	Password: os.Getenv("REDIS_PASSWORD"),
	Protocol: int(protocol),
})

var ctx = context.Background()
var lemmatizer, _ = golem.New(en.New())
var stopwordsMap map[string]struct{}
var nonAlphaRegex, _ = regexp.Compile(`[^a-zA-Z]+`)

func init() {
	// Reading and/or caching of stop words
	result, err := rediscl.SMembers(ctx, "stop_words").Result()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if len(result) == 0 {
		content, err := os.ReadFile("data/stopwords.json")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = json.Unmarshal(content, &result)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, word := range result {
			err := rediscl.SAdd(ctx, "stop_words", word)
			if err != nil {
				fmt.Println(err)
				return
			}
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
	// GETs the page
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request obtained %d status code", resp.StatusCode)
	}

	// Fetchs page content
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Remove all non-text content
	doc.Find("script, style, nav, footer, head, meta, link, noscript").Remove()

	tokenizedContent := strings.Fields(doc.Text())

	// For each token, process and caches it

	terms := make(map[string]struct{})

	for _, word := range tokenizedContent {
		processedWord := ProcessWord(word)
		if processedWord != "" {
			terms[processedWord] = struct{}{}
		}
	}

	pipe := rediscl.Pipeline()

	for term := range terms {
		pipe.SAdd(ctx, term, url)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func FetchUrls(query string) ([]string, error) {
	urls := []string{}

	tokenizedQuery := strings.Fields(query)
	for _, token := range tokenizedQuery {
		currentUrls, err := rediscl.SMembers(ctx, ProcessWord(token)).Result()
		if err != nil {
			fmt.Println(err.Error())
			return []string{}, err
		}

		if len(currentUrls) == 0 {
			continue
		}

		for _, url := range currentUrls {
			if !slices.Contains(urls, url) {
				urls = append(urls, url)
			}
		}
	}

	return urls, nil
}
