package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Kane-Sendgrid/minefield/fs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	fmt.Println("mounting", flag.Arg(0))

	nfs := pathfs.NewPathNodeFs(
		&fs.Fs{FileSystem: pathfs.NewDefaultFileSystem(),
			Files: map[string]*fs.File{}}, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	server.SetDebug(true)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	go func() {
		for sig := range signalChannel {
			fmt.Println("Got signal:", sig, ", attempting graceful shutdown...")
			err := server.Unmount()
			if err != nil {
				fmt.Println("can't unmount due to", err)
			} else {
				return
			}
		}
	}()

	server.Serve()
	println("stopped...")
}
