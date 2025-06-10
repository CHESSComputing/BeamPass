package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	flag.StringVar(&config, "config", "", "server config file")
	flag.Parse()
	if version {
		fmt.Println("server version:", Info())
		return
	}

	conf, err := parseConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Database connection string (replace with your actual credentials)
	// Example: "user:password@tcp(127.0.0.1:3306)/database_name"
	initDB(conf.DBUri)
	defer db.Close()

	// Register the handler for the /btr endpoint
	http.HandleFunc("/btr", BtrHandler)

	// Start the web server
	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
