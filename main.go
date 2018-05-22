package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/machinebox/sdk-go/tagbox"
)

func main() {
	var (
		addr  = flag.String("addr", ":9000", "address")
		state = flag.String("state", "./reddit.machinebox.tagbox", "tagbox state file")
		csv   = flag.String("csv", "./reddit.csv", "cvs file")
	)

	flag.Parse()
	tagbox := tagbox.New("http://localhost:8080")
	fmt.Println(`visualsearch by Machine Box - https://machinebox.io/`)
	fmt.Println("Loading data from csv")
	items, err := LoadData(*csv)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("Done!")

	fmt.Println("Waiting for Tagbox to be ready...")
	boxutil.WaitForReady(context.Background(), tagbox)
	fmt.Println("Done!")

	fmt.Println("Setup tagbox state")

	f, err := os.Open(*state)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = tagbox.PostState(f)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("Done!")
	f.Close()

	fmt.Println("Go to:", *addr+"...")

	srv := NewServer("./assets", tagbox, items)
	if err := http.ListenAndServe(*addr, srv); err != nil {
		log.Fatalln(err)
	}
}
