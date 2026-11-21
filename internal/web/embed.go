package web

import "embed"

//go:embed index.html stats.html preview.html zaps-mark.svg
var Assets embed.FS

