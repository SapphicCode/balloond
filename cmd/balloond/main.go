package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/Pandentia/balloond/balloon"
)

func main() {
	unix := flag.String("unix", "", "Unix socket path to libvirt socket.")
	tcp := flag.String("tcp", "", "TCP address to libvirt socket.")

	interval := flag.Duration("interval", time.Duration(0), "Interval at which to run the ballooning daemon.")
	freeAllowance := flag.Uint64("freeAllowance", 0, "The amount of memory (in kB) to ensure free within each VM.")
	chunkSize := flag.Uint64("memoryChunkSize", 0, "The granularity with which memory is added or removed, in kB.")

	dry := flag.Bool("dry", false, "This flag enables dry run / pretend mode.")
	verbose := flag.Bool("verbose", false, "This flag enables the logger's debug mode, and produces far more output.")

	flag.Parse()

	var conn net.Conn
	var err error
	timeout := 15 * time.Second

	if *unix != "" {
		conn, err = net.DialTimeout("unix", *unix, timeout)
	} else if *tcp != "" {
		conn, err = net.DialTimeout("tcp", *tcp, timeout)
	} else {
		fmt.Println("Either unix or tcp address must be specified.")
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error occurred dialing connection: %s", err)
	}

	b := balloon.New(conn)
	if *verbose {
		b.Logger = b.Logger.Level(zerolog.DebugLevel)
	} else {
		b.Logger = b.Logger.Level(zerolog.InfoLevel)
	}
	if *interval != time.Duration(0) {
		b.Interval = *interval
	}
	if *freeAllowance != 0 {
		b.FreeAllowance = *freeAllowance
	}
	if *chunkSize != 0 {
		b.MemoryChunk = *chunkSize
	}
	b.DryRun = *dry
	b.RunDaemon()
}
