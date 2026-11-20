package web

import "embed"

//go:embed index.html stats.html zaps-mark.svg
var Assets embed.FS

