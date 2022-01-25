package ukuleleweb

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/peterbourgon/diskv/v3"
)

//go:embed templates/*
var templateFiles embed.FS

var baseTmpl = template.Must(template.New("layout").ParseFS(templateFiles, "templates/base/*.html"))
var pageTmpl = template.Must(template.Must(baseTmpl.Clone()).ParseFS(templateFiles, "templates/contents/page.html"))
var editTmpl = template.Must(template.Must(baseTmpl.Clone()).ParseFS(templateFiles, "templates/contents/edit.html"))

type pageValues struct {
	Title         string
	PageName      string
	HTMLContent   template.HTML
	SourceContent string
	Error         string
}

type PageHandler struct {
	MainPage string
	D        *diskv.Diskv
}

func (h *PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/"+h.MainPage, http.StatusMovedPermanently)
		return
	}
	pageName := getPageName(r.URL.Path)
	if pageName == "" {
		http.Error(w, "Invalid page name", http.StatusNotFound)
		return
	}

	tmpl := pageTmpl
	pv := &pageValues{
		Title:    pageName, // XXX insert spaces before capitals
		PageName: pageName,
	}

	if r.FormValue("edit") == "1" {
		tmpl = editTmpl
		content := contentValue(r)
		if content == "" {
			content = h.D.ReadString(pageName)
		}
		pv.SourceContent = content
	} else {
		tmpl = pageTmpl
		if r.Method == "POST" {
			content := contentValue(r)
			err := h.D.WriteString(pageName, content)
			if err != nil {
				pv.Error = err.Error() // XXX hide?
			}
		}

		content := h.D.ReadString(pageName)
		pv.HTMLContent = template.HTML(renderHTML(content))
	}
	err := tmpl.Execute(w, pv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPageName(path string) string {
	if !strings.HasPrefix(path, "/") {
		return ""
	}
	path = path[1:]

	if !isPageName(path) {
		return ""
	}
	return path
}

// isPageName returns true iff pn is a camel case page name.
func isPageName(pn string) bool {
	wantUpper := true
	for _, r := range pn {
		if wantUpper {
			if !(r >= 'A' && r <= 'Z') {
				return false
			}
			wantUpper = false
		}

		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false
		}
	}
	return true
}

func contentValue(r *http.Request) string {
	return strings.ReplaceAll(r.FormValue("content"), "\r\n", "\n")
}

var wikiLinkRE = regexp.MustCompile(`\b([A-Z][a-z]+){2,}\b`)

func renderHTML(md string) string {
	// XXX: It is a hack to replace wiki links after markdown rendering...
	// XXX: Does blackfriday handle wiki links better?
	// return string(blackfriday.MarkdownCommon([]byte(md)))
	md = wikiLinkRE.ReplaceAllStringFunc(md, func(m string) string {
		return fmt.Sprintf(`<a href="/%s">%s</a>`, m, m)
	})

	doc := markdown.Parse([]byte(md), nil)
	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags, // XXX rethink
	})
	return string(markdown.Render(doc, renderer))
}
