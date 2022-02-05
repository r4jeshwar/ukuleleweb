package ukuleleweb

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/peterbourgon/diskv/v3"
)

//go:embed static/*
var StaticFiles embed.FS

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
	ReverseLinks  []string
}

type PageHandler struct {
	MainPage string
	D        *diskv.Diskv

	// A cached version of the reverse links.
	revLinksMu sync.RWMutex
	revLinks   map[string][]string
}

func (h *PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
	w.Header().Set("Referrer-Policy", "no-referrer")

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
		Title:    pageName, // XXX insert spaces before capitals?
		PageName: pageName,
	}

	if r.FormValue("edit") == "1" {
		tmpl = editTmpl
		content := contentValue(r)
		if content == "" {
			content = h.D.ReadString(pageName)
		}
		pv.SourceContent = content
	} else if r.Method == "POST" {
		content := contentValue(r)
		err := h.D.WriteString(pageName, content)
		if err == nil { // Success saving! This is the default case.
			// TODO: Potentially do it in a background job?
			h.recalculateRevLinks()
			http.Redirect(w, r, "/"+pageName, http.StatusFound)
			return
		}
		// On error, render edit form with the error message.
		w.WriteHeader(http.StatusInternalServerError)
		tmpl = editTmpl
		log.Printf("ERROR: diskv.WriteString(%q, ...): %v\n", pageName, err)
		pv.Error = "Internal error writing page"
		pv.SourceContent = content
	} else {
		tmpl = pageTmpl
		content := h.D.ReadString(pageName)
		pv.HTMLContent = template.HTML(renderHTML(content))
		pv.ReverseLinks = h.reverseLinks(pageName)

	}
	err := tmpl.Execute(w, pv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) reverseLinks(pagename string) []string {
	h.revLinksMu.RLock()
	defer h.revLinksMu.RUnlock()

	if h.revLinks == nil {
		// Recalculate at read only if not done before.
		// This should only happen on startup.
		h.revLinksMu.RUnlock()
		h.recalculateRevLinks()
		h.revLinksMu.RLock()
	}
	return h.revLinks[pagename]
}

func (h *PageHandler) recalculateRevLinks() {
	rl := AllReverseLinks(h.D)

	h.revLinksMu.Lock()
	defer h.revLinksMu.Unlock()
	h.revLinks = rl
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
	return pageNameRE.MatchString(pn)
}

func contentValue(r *http.Request) string {
	return strings.ReplaceAll(r.FormValue("content"), "\r\n", "\n")
}
