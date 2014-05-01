package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var configFilename = flag.String("config", "", "pathname of YAML configuration file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: go-pronto -config=/path/to/config")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if *configFilename == "" {
		log.Fatal("Must provide a configuration file")
	}
	conf, err := ConfigFromYaml(*configFilename)
	if err != nil {
		log.Fatalf("Error parsing config file %s", err)
	}
	creds := CredentialsFromConfig(conf)
	pronto := New(creds, conf.Container)

	fmt.Printf("Serving container %s at %s\n", conf.Container, conf.Port)
	log.Fatal(http.ListenAndServe(conf.Port, pronto))
}
