package ui

import "embed"

// comment directive to store ui/html and ui/static as embedded filesystems

//go:embed "html" "static"
var Files embed.FS
