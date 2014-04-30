package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"launchpad.net/goose/client"
	"launchpad.net/goose/identity"
	"launchpad.net/goose/swift"
)

var configFilename = flag.String("config", "", "pathname of YAML configuration file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: go-pronto -config=/path/to/config")
		flag.PrintDefaults()
	}
}

type ProntoService struct {
	s         *swift.Client
	container string
}

func New(creds *identity.Credentials, container string) *ProntoService {
	cl := client.NewClient(creds, identity.AuthUserPass, nil)
	sw := swift.New(cl)
	return &ProntoService{sw, container}
}

func (p ProntoService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.TrimLeft(r.URL.Path, "/")
	log.Printf("Request for %#v\n", trimmedPath)
	rc, err := p.s.GetReader(p.container, trimmedPath)
	if err != nil {
		log.Printf("%s", err)
		http.HandlerFunc(http.NotFound)(w, r)
		return
	}
	defer rc.Close()
	io.Copy(w, rc)
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
