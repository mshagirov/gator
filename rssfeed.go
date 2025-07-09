package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPError: Status Code %v", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rssData RSSFeed
	err = xml.Unmarshal(bodyBytes, &rssData)
	if err != nil {
		return nil, err
	}
	rssData.Channel.Title = html.UnescapeString(rssData.Channel.Title)
	rssData.Channel.Description = html.UnescapeString(rssData.Channel.Description)
	for i := range rssData.Channel.Item {
		rssData.Channel.Item[i].Title = html.UnescapeString(rssData.Channel.Item[i].Title)
		rssData.Channel.Item[i].Description = html.UnescapeString(rssData.Channel.Item[i].Description)
	}
	return &rssData, nil
}
