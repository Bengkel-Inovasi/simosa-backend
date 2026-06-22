package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type NewsArticle struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Source  string `json:"source"`
	PubDate string `json:"pub_date"`
}

func GetEconomicNews(w http.ResponseWriter, r *http.Request) {
	// CNBC Indonesia Market RSS
	rssURL := "https://www.cnbcindonesia.com/market/rss"

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(rssURL)
	if err != nil {
		fmt.Printf("Error fetching RSS feed: %v, using fallback news\n", err)
		serveFallbackNews(w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("RSS feed returned status code %d, using fallback news\n", resp.StatusCode)
		serveFallbackNews(w)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading RSS feed body: %v, using fallback news\n", err)
		serveFallbackNews(w)
		return
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Printf("Error unmarshaling RSS feed XML: %v, using fallback news\n", err)
		serveFallbackNews(w)
		return
	}

	articles := make([]NewsArticle, 0)
	limit := 6
	if len(rss.Channel.Items) < limit {
		limit = len(rss.Channel.Items)
	}

	for i := 0; i < limit; i++ {
		item := rss.Channel.Items[i]
		articles = append(articles, NewsArticle{
			Title:   item.Title,
			Link:    item.Link,
			Source:  "CNBC Indonesia",
			PubDate: item.PubDate,
		})
	}

	if len(articles) == 0 {
		serveFallbackNews(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func serveFallbackNews(w http.ResponseWriter) {
	fallback := []NewsArticle{
		{
			Title:   "Harga CPO Naik Menjadi Rp 13.200/kg Akibat Kenaikan Permintaan Global",
			Link:    "https://www.cnbcindonesia.com",
			Source:  "Info Sawit",
			PubDate: time.Now().Add(-1 * time.Hour).Format(time.RFC1123),
		},
		{
			Title:   "Ekspor Minyak Sawit RI ke India Menguat di Kuartal II 2026",
			Link:    "https://www.cnbcindonesia.com",
			Source:  "Warta Ekonomi",
			PubDate: time.Now().Add(-4 * time.Hour).Format(time.RFC1123),
		},
		{
			Title:   "Kebijakan B40 Dorong Konsumsi Domestik Minyak Kelapa Sawit",
			Link:    "https://www.cnbcindonesia.com",
			Source:  "Kementerian ESDM",
			PubDate: time.Now().Add(-12 * time.Hour).Format(time.RFC1123),
		},
		{
			Title:   "Petani Sawit Mandiri Didorong Terapkan Praktik Berkelanjutan ISPO",
			Link:    "https://www.cnbcindonesia.com",
			Source:  "Ditjen Perkebunan",
			PubDate: time.Now().Add(-24 * time.Hour).Format(time.RFC1123),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallback)
}
