package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func Info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var config string
	cfile := fmt.Sprintf("%s/.beampass.json", os.Getenv("HOME"))
	flag.StringVar(&config, "config", cfile, "server config file")
	flag.Parse()
	if version {
		fmt.Println("server version:", Info())
		return
	}

	conf, err := parseConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// initialize connection to beampass user db with dbURI, e.g.
	// user:password@tcp(127.0.0.1:3306)/database_name
	initDB(conf.DBUri)
	defer db.Close()

	// Register the handler for the /btr endpoint
	http.HandleFunc("/btr", BtrHandler)

	// Start the web server
	port := fmt.Sprintf(":%d", conf.Port)
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
