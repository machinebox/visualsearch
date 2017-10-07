package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/machinebox/sdk-go/tagbox"
)

func main() {
	tbox := tagbox.New("http://localhost:8080")
	f, err := os.Open("../reddit.csv")
	if err != nil {
		log.Fatal(err)
		return
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
			log.Fatal(err)
		}

		u, err := url.Parse(record[1])
		if err != nil {
			continue
		}
		err = tbox.TeachURL(u, record[0], "other")
		if err != nil {
			log.Printf("[ERROR] teaching %v -> %v %v\n", record[0], record[1], err)
		} else {
			log.Printf("[INFO] teach %v\n", record[0])
		}
	}

}
