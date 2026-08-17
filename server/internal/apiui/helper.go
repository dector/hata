package apiui

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
)

const tokenStorageKey = "hata_apiui_token"

type pageViewData struct {
	Title       string
	Heading     string
	Description template.HTML
	Path        string
	BackHref    string
	Content     template.HTML
	CurlTitle   string
	CurlSnippet string
	Script      template.JS
}

//go:embed page.html
var pageHTML string

var pageTemplate = template.Must(template.New("apiui-page").Parse(pageHTML))

func writePage(w http.ResponseWriter, data pageViewData) {
	type pageTemplateData struct {
		Title       string
		Heading     string
		Description template.HTML
		Path        string
		BackHref    string
		Content     template.HTML
		CurlTitle   string
		CurlSnippet string
		Script      template.JS
		TokenKey    string
	}

	tplData := pageTemplateData{
		Title:       data.Title,
		Heading:     data.Heading,
		Description: data.Description,
		Path:        data.Path,
		BackHref:    data.BackHref,
		Content:     data.Content,
		CurlTitle:   data.CurlTitle,
		CurlSnippet: data.CurlSnippet,
		Script:      data.Script,
		TokenKey:    tokenStorageKey,
	}

	var buf bytes.Buffer
	if err := pageTemplate.Execute(&buf, tplData); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}
