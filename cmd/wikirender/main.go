package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gnoack/ukuleleweb"
)

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		md, err := os.ReadFile(fn)
		if err != nil {
			log.Fatalf("ReadFile(%q): %v\n", fn, err)
			continue
		}
		html, err := ukuleleweb.RenderHTML(string(md))
		if err != nil {
			log.Fatalf("RenderHTML(ReadFile(%q)): %v", fn, err)
			continue
		}
		fmt.Print(html)
	}
}
