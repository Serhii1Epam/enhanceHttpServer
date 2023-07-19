/*
 * Cmd start Enhance HTTP server
 */
package main

import (
	"fmt"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appserver"
)

func main() {
	fmt.Println("Enhance HTTP Server start...")
	defer fmt.Println("Enhance HTTP Server finish work.")
	srv := appserver.NewServer()
	if !srv.IsRun() {
		srv.SrvRun()
	}
}
