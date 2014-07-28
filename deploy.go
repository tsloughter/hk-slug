package main

import (
	"fmt"
	"github.com/cyberdelia/heroku-go/v3"
	"os"
	"path/filepath"
)

const (
	PLUGIN_NAME    = "deploy"
	PLUGIN_VERSION = 1
	ENDPOINT       = "https://hk-deploy.herokuapp.com/slot"
	INFO_PREAMBLE  = `%s %d: Deploy code to Heroku using the API`
	HELP_TEXT      = `Usage: hk deploy DIRECTORY

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

	dir := os.Args[1] // TODO: Maybe fallback to CWD or Git root?

	fullPath, _ := filepath.Abs(dir)
	fmt.Printf("Creating .tgz of %s...\n", fullPath)
	tgz := buildTgz(dir)
	fmt.Printf("done (%d bytes)\n", tgz.Len())

	fmt.Print("Requesting upload slot... ")
	slot, err := getUploadSlot()
	if err == nil {
		fmt.Println("done")
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print("Uploading .tgz to S3... ")
	if err := upload(&tgz, slot); err == nil {
		fmt.Println("done")
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// working downloadable link now in slot.DownloadUrl
	fmt.Print("Submitting build with download link... ")
	if _, err := submitBuild(&slot.DownloadUrl); err == nil {
		fmt.Println("done")
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// TODO: stream output (https://devcenter.heroku.com/articles/build-and-release-using-the-api#experimental-realtime-build-output).
	//       To do so, the heroku client will need to be updated to work with `edge` schema.

}

func submitBuild(url *string) (*heroku.Build, error) {
	app := os.Getenv("HKAPP")
	heroku.DefaultTransport.Username = os.Getenv("HKUSER")
	heroku.DefaultTransport.Password = os.Getenv("HKPASS")

	hk := heroku.NewService(heroku.DefaultClient)

	o := new(heroku.BuildCreateOpts)
	o.SourceBlob.URL = url
	// TODO: allow specifiying o.Version to a custom value and/or inferring it

	if build, err := hk.BuildCreate(app, *o); err != nil {
		return nil, err
	} else {
		return build, nil
	}
}
