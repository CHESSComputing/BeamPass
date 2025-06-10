package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	DBUri string
	Port  int
}

func parseConfig(fname string) (Config, error) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var c Config
	err = json.Unmarshal(data, &c)
	return c, err
}
