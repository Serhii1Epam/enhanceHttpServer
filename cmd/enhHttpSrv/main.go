/*
 * Cmd start Enhance HTTP server
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appserver"
)

var cpuprofile = flag.String("cpuprofile", "testcpu.out", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "testmem.out", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	fmt.Println("Enhance HTTP Server start...")
	defer fmt.Println("Enhance HTTP Server finish work.")
	srv := appserver.NewServer()
	if !srv.IsRun() {
		srv.SrvRun()
	}
}
