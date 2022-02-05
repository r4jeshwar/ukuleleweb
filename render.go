package ukuleleweb

import (
	"regexp"
	"sort"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/peterbourgon/diskv/v3"
)

func renderHTML(md string) string {
	// XXX: Does blackfriday handle wiki links better?
	// return string(blackfriday.MarkdownCommon([]byte(md)))

	// XXX: It is a hack to replace wiki links before markdown rendering...
	md = replaceDetectedLinks(md)

	doc := markdown.Parse([]byte(md), nil)
	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags, // XXX rethink
	})
	return string(markdown.Render(doc, renderer))
}

var (
	pageNameRE = regexp.MustCompile(`\b([A-ZÄÖÜ][a-zäöüß]+){2,}\b`)
	goLinkRE   = regexp.MustCompile(`\bgo/[A-Za-z0-9_+öäüÖÄÜß-]+\b`)

	pageNameWithPrefixRE = regexp.MustCompile(`(\s)(` + pageNameRE.String() + `)`)
)

func replaceDetectedLinks(t string) string {
	t = pageNameWithPrefixRE.ReplaceAllString(" "+t, `$1<a href="/$2">$2</a>`)[1:]
	t = goLinkRE.ReplaceAllString(t, `<a href="http://$0">$0</a>`)
	return t
}

// OutgoingLinks returns the outgoing wiki links in a given Markdown input.
// The outgoing links are a map of page names to true.
func OutgoingLinks(md string) map[string]bool {
	res := make(map[string]bool)
	for _, m := range pageNameWithPrefixRE.FindAllStringSubmatch(" "+md, -1) {
		res[m[2]] = true
	}
	return res
}

// ReverseLinks calculates the reverse link map for the whole wiki.
// The returned map maps page names to a list of pages linking to it.
// Sets of pages are represented as sorted lists.
func AllReverseLinks(d *diskv.Diskv) map[string][]string {
	revLinks := make(map[string]map[string]bool)
	for p := range d.Keys(nil) {
		pOut := OutgoingLinks(d.ReadString(p))
		for q, _ := range pOut {
			qIn, ok := revLinks[q]
			if !ok {
				qIn = make(map[string]bool)
				revLinks[q] = qIn
			}
			qIn[p] = true
		}
	}

	revLinksSorted := make(map[string][]string)
	for p, s := range revLinks {
		revLinksSorted[p] = sortedStringSlice(s)
	}
	return revLinksSorted
}

func sortedStringSlice(a map[string]bool) []string {
	var res []string
	for k, _ := range a {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}
