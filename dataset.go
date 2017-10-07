package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Item struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`

	Confidence float64 `json:"confidence,omitempty"`
}

func LoadData(path string) (map[string]Item, error) {
	items := map[string]Item{}
	f, err := os.Open(path)
	if err != nil {
		return items, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		// 0:id, 1:url, 2:title
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return items, err
		}
		items[record[0]] = Item{
			ID:    record[0],
			URL:   record[1],
			Title: strings.TrimSpace(record[2]),
		}
	}
	return items, nil
}
