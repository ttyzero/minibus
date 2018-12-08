package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ttyzero/minibus/lib"
)

func main() {
	// Find the cache directory for our working dir
	workDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal("Could not find user cache directory", err)
	}
	workDir = filepath.Join(workDir, "minibus")

	// Disable log output unless DEBUG is set
	if os.Getenv("DEBUG") == "" {
		log.SetOutput(ioutil.Discard)
	}

	// Start the bus!
	lib.Start(workDir)
}
