package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"os"
	"log"

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

	var err error
	tarball := os.Args[1] // TODO: Maybe fallback to CWD or Git root?

	app := os.Getenv("HKAPP")
	apiKey := os.Getenv("HKPASS")

	fmt.Println("Deploying... ")
	dir := unpack(tarball)
	s := slug.NewSlug(apiKey, app, dir)
	tarFile := s.Archive()
	fmt.Printf("done\nPushing %s... ", tarFile.Name())
	err = s.Push()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("done\n")
	fmt.Printf("Releasing... ")
	release := s.Release()
	fmt.Printf("done (v%d)\n", release.Version)

	fmt.Printf("Completed deploy of v%d\n", release.Version)
}

func unpack(tarball string) string {
	gzipFile, _ := os.Open(tarball)
	tarFile, _ := gzip.NewReader(gzipFile)
	defer gzipFile.Close()

	tr := tar.NewReader(tarFile)
	slugDir, _ := ioutil.TempDir("", "slug")
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		newDir := slugDir + "/" + filepath.Dir(hdr.Name)
		err = os.MkdirAll(newDir, os.ModeDir | 0766)
		if err != nil {
			log.Fatalln(err)
		}

		newFile := newDir + "/" + filepath.Base(hdr.Name)
		fo, err := os.Create(newFile)
		if err != nil {
			log.Fatalln(err)
		}
		fo.Chmod(hdr.FileInfo().Mode())
		if _, err := io.Copy(fo, tr); err != nil {
			log.Fatalln(err)
		}
	}

	return slugDir
}
