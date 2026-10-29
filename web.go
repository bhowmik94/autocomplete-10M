package main

import "embed"

//go:embed web/index.html
var webFiles embed.FS
