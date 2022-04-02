// An experimental command line tool for visualizing the graph of wiki pages.
// The output of this command is a digraph for input with GraphViz tools.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gnoack/ukuleleweb"
	"github.com/peterbourgon/diskv/v3"
)

var storeDir = flag.String("store_dir", "", "Store directory")

func writeDigraph(w io.Writer, d *diskv.Diskv) {
	fmt.Fprintln(w, "digraph G {")
	fmt.Fprintln(w, "\toverlap = false;")
	fmt.Fprintln(w, "\tnode [color=red];")

	fmt.Fprintln(w)
	for pn := range d.Keys(nil) {
		fmt.Fprintf(w, "\t%v [color=black shape=box];\n", pn)
	}

	fmt.Fprintln(w)
	for pn := range d.Keys(nil) {
		md := d.ReadString(pn)
		for ogPn, _ := range ukuleleweb.OutgoingLinks(md) {
			fmt.Fprintf(w, "\t%v -> %v;\n", pn, ogPn)
		}
	}
	fmt.Fprintln(w, "}")
}

func main() {
	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(o, "\t%s -store_dir=/path/to/wiki | neato -Tsvg > out.svg\n", os.Args[0])
		fmt.Fprintln(o)
		fmt.Fprintln(o, "Flags:")
		flag.PrintDefaults()
	}

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
	writeDigraph(os.Stdout, d)
}
