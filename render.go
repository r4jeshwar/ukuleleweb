package ukuleleweb

import (
	"fmt"
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

var goLinkRE = regexp.MustCompile(`\bgo/[A-Za-z0-9_+öäüÖÄÜß-]+\b`)

func renderHTML(md string) string {
	// XXX: Does blackfriday handle wiki links better?
	// return string(blackfriday.MarkdownCommon([]byte(md)))

	// XXX: It is a hack to replace wiki links before markdown rendering...
	md = pageNameRE.ReplaceAllStringFunc(md, func(m string) string {
		return fmt.Sprintf(`<a href="/%s">%s</a>`, m, m)
	})
	// Go links.
	md = goLinkRE.ReplaceAllStringFunc(md, func(m string) string {
		return fmt.Sprintf(`<a href="http://%s">%s</a>`, m, m)
	})

	doc := markdown.Parse([]byte(md), nil)
	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags, // XXX rethink
	})
	return string(markdown.Render(doc, renderer))
}
