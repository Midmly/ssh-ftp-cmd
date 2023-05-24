package utils

import (
	"log"
	"sync"
	"testing"
)

func TestNewSsl(t *testing.T) {
	cli := &SshConfig{
		Username:   "yeastar",
		Password:   "t1Max6KdD24i5zPvBNnf",
		Address:    "213.9.29.138",
		Port:       1022,
		PrivateKey: "",
		Width:      0,
		Height:     0,
	}
	err := cli.NewSshClient()
	if err != nil {
		log.Print(err)
		return
	}
	var g sync.WaitGroup
	g.Add(1)
	go func() {
		str1 := cli.RunCmd("ls")
		log.Printf("str1: %s", str1)
		g.Done()
	}()
	g.Add(1)
	go func() {
		str2 := cli.RunCmd("echo 't1Max6KdD24i5zPvBNnf'|sudo su")
		log.Printf("str2: %s", str2)
		g.Done()
	}()
	g.Wait()
}
