package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	err := app.f.Split(cleanPath(src), cleanPath(dst), ovr, flat)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func cleanPath(p string) string {
	if len(p) > 0 && p[len(p)-1] == '/' {
		return p[:len(p)-1]
	}

	absPath, err := resolvePath(p)
	if err != nil {
		panic(fmt.Sprintf("Error resolving path %s: %v\n", p, err))

	}

	return absPath
}

func resolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	if path[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func init() {
	flag.StringVar(&src, "src", "", "Source file path")
	flag.StringVar(&dst, "dst", "", "Destination file path")
	flag.BoolVar(&ovr, "overwrite", false, "--overwrite")
	flag.BoolVar(&flat, "flat", false, "--flat")
	flag.Parse()
}
