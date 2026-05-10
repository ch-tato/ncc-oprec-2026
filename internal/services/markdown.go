package services

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// md is the pre-configured Goldmark instance.
var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavoured Markdown (tables, strikethrough, etc.)
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			// Raw HTML in Markdown source is escaped by default (safe).
		),
	)
}

// RenderMarkdown converts raw Markdown text to sanitised HTML.
// Unsafe raw HTML tags in the source are stripped by Goldmark.
func RenderMarkdown(source string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
