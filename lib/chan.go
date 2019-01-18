package lib

import (
	"log"
	"net"
)

type Channel struct {
	name  string
	stop  chan bool
	Send  chan []byte
	socks []net.Conn
}

func NewChannel(name string, busStop chan bool) *Channel {
	log.Printf("Creating channel '%s'", name)

	stop := make(chan bool)
	send := make(chan []byte)
	socks := make([]net.Conn, 0)
	c := Channel{name, stop, send, socks}
	go c.accept()
	return &c
}

func (t *Channel) accept() {
	for {
		select {
		case d := <-t.Send:
			for _, sock := range t.socks {
				_, err := sock.Write(d)
				if err != nil {
					log.Fatal("write error:", err)
				}
			}
		case <-t.stop:
			log.Printf("Shutting down channel, closing sockets")
			for _, sock := range t.socks {
				sock.Close()
				// TODO remove from socks if it closes
			}
		}
	}

}

func (t *Channel) Connect(sockPath string) {
	log.Printf("Connecting to channel %s", sockPath)

	c, err := net.Dial("unixgram", sockPath)
	if err != nil {
		panic(err)
	}
	t.socks = append(t.socks, c)
}
