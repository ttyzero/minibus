package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide a channel name")
		os.Exit(1)
	}
	channel := os.Args[1]
	fmt.Println("Connecting to channel:", channel)
	workDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	pid := []byte(strconv.Itoa(os.Getpid()))

	sockPath := filepath.Join(workDir, "minibus", fmt.Sprintf("%s-%s", pid, channel))
	fmt.Println(sockPath)

	addr, err := net.ResolveUnixAddr("unixgram", sockPath)
	if err != nil {
		fmt.Println("Could not resolve unix socket: " + err.Error())
		os.Exit(1)
	}

	c, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		fmt.Println("what" + err.Error())
		os.Exit(1)
	}

	defer func() {
		c.Close()
		os.Remove(sockPath)
	}()

	for {
		var buf [1024]byte
		n, err := c.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(buf[:n]))
	}

}
