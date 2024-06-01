package main

import (
	"flag"
	"fmt"
)

var (
	app  = application{f: &files{}}
	src  string
	dst  string
	ovr  bool
	flat bool
)

type application struct {
	f filer
}

func main() {
	if src == "" || dst == "" {
		flag.Usage()
		fmt.Println("Both src and dst flags are required")
		return
	}

	app.f.Split(cleanPath(src), cleanPath(dst), ovr, flat)
}

func cleanPath(p string) string {
	if len(p) > 0 && p[len(p)-1] == '/' {
		return p[:len(p)-1]
	}

	return p
}

func init() {
	flag.StringVar(&src, "src", "", "Source file path")
	flag.StringVar(&dst, "dst", "", "Destination file path")
	flag.BoolVar(&ovr, "overwrite", false, "--overwrite")
	flag.BoolVar(&flat, "flat", false, "--flat")
	flag.Parse()
}
