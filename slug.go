package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/naaman/slug"
)

const (
	PLUGIN_NAME    = "slug"
	PLUGIN_VERSION = 1
	ENDPOINT       = "https://hk-deploy.herokuapp.com/slot"
	INFO_PREAMBLE  = `%s %d: Deploy pre-built code to Heroku using the API`
	HELP_TEXT      = `Usage: hk slug TARBALL

	Deploy the specified directory to Heroku`
)

func help() {
	fmt.Println(HELP_TEXT)
}

func info() {
	fmt.Printf(INFO_PREAMBLE+"\n\n", PLUGIN_NAME, PLUGIN_VERSION)
	help()
	os.Exit(0)
}

func main() {
	if os.Getenv("HKPLUGINMODE") == "info" {
		info()
	}

	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
		os.Exit(0)
	}

	tarball := os.Args[1] // TODO: Maybe fallback to CWD or Git root?

	fullPath, _ := filepath.Abs(tarball)
	tar, _ := os.Open(fullPath)
	app := os.Getenv("HKAPP")
	apiKey := os.Getenv("HKPASS")

	fmt.Println("Deploying...")
	myslug := slug.NewSlug(apiKey, app, fullPath)
	myslug.SetArchive(tar)
	myslug.Push()
	release := myslug.Release()
	fmt.Printf("Done (v%d)\n", release.Version)
}
