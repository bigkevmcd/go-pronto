package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"launchpad.net/goose/client"
	"launchpad.net/goose/identity"
	"launchpad.net/goose/swift"
)

type swiftReader interface {
	GetReader(containerName, objectName string) (io.ReadCloser, http.Header, error)
}

type ProntoService struct {
	s         swiftReader
	container string
}

func New(creds *identity.Credentials, container string) *ProntoService {
	cl := client.NewClient(creds, identity.AuthUserPass, nil)
	sw := swift.New(cl)
	return &ProntoService{sw, container}
}

func (p *ProntoService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.TrimLeft(r.URL.Path, "/")
	log.Printf("Request for %#v\n", trimmedPath)
	rc, headers, err := p.s.GetReader(p.container, trimmedPath)
	if err != nil {
		log.Printf("%s", err)
		http.HandlerFunc(http.NotFound)(w, r)
		return
	}

	// Transfer the headers, so we get the correct content-type and Etags etc.
	for k, v := range headers {
		w.Header().Set(k, v[0])
	}
	defer rc.Close()
	io.Copy(w, rc)
}
