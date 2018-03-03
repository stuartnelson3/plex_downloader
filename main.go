package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type request struct {
	Destination string `json:"destination"`
	Link        string `json:"link"`
}

func (r *request) path() string {
	return strings.Split(r.Link, *pathSplit)[1]
}

func (r *request) dst() string {
	return "/var/lib/plexmediaserver/" + r.Destination
}

var (
	srcServer = flag.String("src.server", "", "server with the data")
	pathSplit = flag.String("split", "", "where to split incoming link data")
)

func main() {
	sftpc := make(chan *request)
	ctx := context.Background()

	flag.Parse()

	go func() {
		for req := range sftpc {
			go func() {
				fmt.Printf("starting request for %s to %s\n", req.path(), req.dst())
				cmd := exec.CommandContext(ctx, "sftp", "-r", fmt.Sprintf("%s:%s", *srcServer, req.path()), req.dst())

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("couldn't run command: %v\n", err)
					return
				}
			}()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, fmt.Sprintf("couldn't decode body: %v\n", err), http.StatusBadRequest)
			return
		}

		sftpc <- req

		w.Write([]byte(fmt.Sprintf("downloading %s", req.Link)))
	})

	http.ListenAndServe(":4567", nil)
}
