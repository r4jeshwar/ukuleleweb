package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gnoack/ukuleleweb"
	"github.com/peterbourgon/diskv/v3"
)

var (
	listenAddr = flag.String("addr", "localhost:8080", "HTTP listen address")
	storeDir   = flag.String("store_dir", "", "Store directory")
	mainPage   = flag.String("main_page", "MainPage", "The default page to use as the main page")
)

func main() {
	flag.Parse()

	if *storeDir == "" {
		fmt.Fprintln(flag.CommandLine.Output(), "Needs --store_dir")
		flag.Usage()
		return
	}

	d := diskv.New(diskv.Options{
		BasePath:     *storeDir,
		CacheSizeMax: 1024 * 1024, // 1MB
	})
	http.Handle("/static/", http.FileServer(http.FS(ukuleleweb.StaticFiles)))
	http.Handle("/", &ukuleleweb.PageHandler{
		MainPage: *mainPage,
		D:        d,
	})

	fmt.Printf("Listening on http://%s/\n", *listenAddr)
	err := http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		log.Printf("http.ListenAndServe: %v", err)
	}
}
