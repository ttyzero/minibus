package lib

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

func Start(workDir string) {
	pidFile := filepath.Join(workDir, "minibus.pid")

	// create our working directory if it doesn't already exist
	log.Printf("Checking '%s' exists", workDir)
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
			log.Printf("Removing stale pidfile '%s'\n", pidFile)
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
	log.Println("Minibus open for bus-iness")
	<-bus.stop
}

type Minibus struct {
	path     string
	in       chan (string)
	stop     chan (bool)
	channels map[string]*Channel
	fsw      *fsnotify.Watcher
}

func NewMinibus(path string) (*Minibus, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	in := make(chan (string))
	stop := make(chan (bool))
	bus := Minibus{path, in, stop, map[string]*Channel{}, w}
	bus.fsw.Add(path)
	go bus.monitor()
	go func() {
		sigint := make(chan os.Signal, 2)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		stop <- true
	}()
	return &bus, nil
}

// bus.monitor monitors the bus working directory and creates/deletes channels
func (bus *Minibus) monitor() {
	for {
		select {
		case event, ok := <-bus.fsw.Events:
			if !ok {
				return
			}

			// catch create channels
			if event.Op&fsnotify.Create == fsnotify.Create {
				if filepath.Dir(event.Name) == bus.path {
					finfo, err := os.Stat(event.Name)
					if err != nil {
						log.Printf("Cannot stat '%s'", event.Name)
						log.Println(err)
					}
					if finfo.IsDir() {
						bus.newChannel(event.Name)
					}
				}
			}

			//catch potential removed channels
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if filepath.Dir(event.Name) == bus.path {
					bus.closeChannel(event.Name)
				}
			}

		case err, ok := <-bus.fsw.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (bus *Minibus) newChannel(path string) {
	ch, err := NewChannel(path, bus.stop)
	if err != nil {
		log.Println("couldn't open channel", err)
		return
	}
	bus.channels[path] = ch
}

func (bus *Minibus) closeChannel(path string) {
	log.Printf("Destroyed channel '%s'", path)

}

type Sock struct {
	pid int
}
