package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

const workDir string = "/tmp/minibus"

func main() {
	pidFile := filepath.Join(workDir, "minibus.pid")

	// create our working directory if it doesn't already exist
	err := os.MkdirAll(workDir, 0766)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create working directory '%s'\n", workDir)
		os.Exit(1)
	}

	// check to see that another bus is not already running this bus route
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		data, _ := ioutil.ReadFile(pidFile)
		oldPid, _ := strconv.Atoi(string(data))
		oldProc, err := os.FindProcess(oldPid)
		if err == nil {
			err := oldProc.Signal(syscall.Signal(0))
			if err == nil {
				fmt.Fprintf(os.Stderr, "Another bus is running '%d'\n", oldPid)
				os.Exit(1)
			}
			fmt.Printf("Removing stale pidfile '%s'\n", pidFile)
			err = os.Remove(pidFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot remove stale pidfile '%d'\n", pidFile)
				os.Exit(1)
			}
		}
	}

	// write a PID file declaring our owndership of this bus route
	err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	defer func() {
		_ = os.Remove(pidFile)
	}()

	bus, err := NewMinibus(workDir)
	if err != nil {
		log.Println(err)
	}
	<-bus.stop
}

type Minibus struct {
	path     string
	in       chan (string)
	stop     chan (bool)
	channels []*Channel
	fsw      *fsnotify.Watcher
}

func NewMinibus(path string) (*Minibus, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	in := make(chan (string))
	stop := make(chan (bool))
	bus := Minibus{path, in, stop, []*Channel{}, w}
	bus.fsw.Add(path)
	// bus.open()
	go bus.monitor()
	go func() {
		sigint := make(chan os.Signal, 2)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		stop <- true
	}()
	return &bus, nil
}

func (bus *Minibus) monitor() {
	for {
		select {
		case event, ok := <-bus.fsw.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
		case err, ok := <-bus.fsw.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

type Channel struct {
	in    chan (string)
	socks []*Sock
}

func NewChannel(path string) *Channel {
	in := make(chan (string))
	c := Channel{in, []*Sock{}}
	return &c
}

type Sock struct {
}
