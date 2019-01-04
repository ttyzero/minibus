package lib

import (
	"log"
)

type Channel struct {
	name  string
	stop  chan bool
	socks map[int]*Sock
}

func NewChannel(name string, busStop chan bool) *Channel {
	log.Printf("Creating channel '%s'", name)

	stop := make(chan bool)
	c := Channel{name, stop, map[int]*Sock{}}
	return &c
}

func (t *Channel) Connect(sockPath string) error {
	log.Printf("Connecting to channel %s", sockPath)
	return nil
}

func (t *Channel) Disconnect(sockPath string) error {

	return nil
}
