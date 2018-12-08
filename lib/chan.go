package lib

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Channel struct {
	path  string
	name  string
	in    chan string
	stop  chan bool
	socks map[string]*Sock
	fsw   *fsnotify.Watcher
}

func NewChannel(path string, busStop chan bool) (*Channel, error) {
	name := filepath.Base(path)
	log.Printf("Creating channel '%s'", name)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	in := make(chan string)
	stop := make(chan bool)
	c := Channel{path, name, in, stop, map[string]*Sock{}, w}
	c.fsw.Add(path)
	go c.monitor(busStop)
	return &c, nil
}

func (ch *Channel) close() {
	close(ch.stop)
}

// Channel.monitor monitors the channel dir and connects/disconnects sockets
func (ch *Channel) monitor(busStop chan bool) {
	for {
		select {
		case <-busStop:
			// TODO cleanup in sock, disconnect other sockets
			ch.close()
			return
		case event, ok := <-ch.fsw.Events:
			if !ok {
				return
			}

			// catch create channels
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Printf("Can '%s'", event.Name)
			}

			//catch potential removed channels
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Printf("Cannot '%s'", event.Name)
			}

		case err, ok := <-ch.fsw.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
