package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
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
		if pidActive(oldPid) {
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

var (
	MSG_RE      = regexp.MustCompile("^(\\w+):(.*)")
	SOCKFILE_RE = regexp.MustCompile("^([0-9]*)-([-_\\.\\w]*)")
)

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
	go bus.datagramListener()
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
				// ignore our own minibus socket
				if filepath.Base(event.Name) != "minibus" {
					bus.add(event.Name)
				}
			}

			//catch potential removed channels
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if filepath.Dir(event.Name) == bus.path {
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

func (bus *Minibus) add(path string) {
	if filepath.Dir(path) == bus.path {
		finfo, err := os.Stat(path)
		if err != nil {
			log.Printf("Cannot stat '%s'", path)
			log.Println(err)
		}

		if finfo.Mode()&os.ModeSocket != 0 {
			m := SOCKFILE_RE.FindStringSubmatch(filepath.Base(path))
			if len(m) == 3 {
				pid := m[1]
				chn := m[2]
				if _, ok := bus.channels[chn]; !ok {
					bus.channels[chn] = NewChannel(path, bus.stop)
				}
				go bus.channels[chn].Connect(path)
				log.Printf("[%s]:%s", chn, pid)
			} else {
				log.Printf("cannot parse filename, should be $PID-$CHANNEL: %s", filepath.Base(path))
			}
		}
	}
}
func (bus *Minibus) datagramListener() {
	sockPath := filepath.Join(bus.path, "minibus")

	err := syscall.Unlink(sockPath)
	if err != nil {
		log.Println("Unlink()", err)
	}

	addr, err := net.ResolveUnixAddr("unixgram", sockPath)
	if err != nil {
		println("Could not resolve unix socket: " + err.Error())
		os.Exit(1)
	}

	// listen on the socket
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		println("Could not listen on unix socket datagram: " + err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		bus.handleMsg(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// close socket when we finish
	defer conn.Close()
}

func (bus *Minibus) handleMsg(data string) {
	var chn, msg string
	log.Println("----------------------")
	m := MSG_RE.FindStringSubmatch(data)
	if len(m) == 3 {
		chn = m[1]
		msg = m[2]
		log.Printf("[%s]:%s", chn, msg)

		if _, ok := bus.channels[chn]; ok {
			bus.channels[chn].Send <- []byte(msg)
		}
	} else {
		log.Printf("cannot parse: %s", data)
	}

}

type Sock struct {
	pid int
}
