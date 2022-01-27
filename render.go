package ukuleleweb

import (
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
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
)

func replaceDetectedLinks(t string) string {
	t = pageNameRE.ReplaceAllString(t, `<a href="/$0">$0</a>`)
	t = goLinkRE.ReplaceAllString(t, `<a href="http://$0">$0</a>`)
	return t
}
