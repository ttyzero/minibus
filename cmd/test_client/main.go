package main

import (
	"fmt"
	"os"

	minibus "github.com/ttyzero/minibus-go"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide a channel name")
		os.Exit(1)
	}
	channel := os.Args[1]

	mb := minibus.New(minibus.Default)

	/*
		err := mb.Send(channel, []byte("Sending"))
		if err != nil {
			fmt.Println(err)
		}
	*/
	chTest, closeTest, err := mb.OpenChannel(channel)
	if err != nil {
		fmt.Println("faul", err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case msg := <-chTest:
				fmt.Println("Got a msg:", msg)
			case <-closeTest:
				fmt.Println("closing here..")
				return
			}
		}
	}()
	<-closeTest
}
