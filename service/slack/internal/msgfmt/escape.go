package msgfmt

import (
	"strings"
)

var unescapeReplacer = strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">")

func Unescape(s string) string {
	return unescapeReplacer.Replace(s)
}

var escapeReplacer = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func Escape(s string) string {
	return escapeReplacer.Replace(s)
}
