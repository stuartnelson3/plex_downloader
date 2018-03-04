package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type request struct {
	Destination string `json:"destination"`
	Link        string `json:"link"`
}

func (r *request) path() string {
	p, err := url.PathUnescape(r.Link)
	if err != nil {
		panic(err)
	}
	return strings.Split(p, *pathSplit)[1]
}

func (r *request) dst() string {
	return "/var/lib/plexmediaserver/" + r.Destination
}

func (r *request) sftpCmd(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, "sftp", "-r", fmt.Sprintf(`%s:"%s"`, *srcServer, r.path()), r.dst())
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
				start := time.Now()
				fmt.Printf("starting request for %s to %s\n", req.path(), req.dst())

				cmd := req.sftpCmd(ctx)

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("couldn't run command: %v\n", err)
					return
				}
				fmt.Printf("finished downloading %s to %s. time: %v\n", req.path(), req.dst(), time.Since(start))
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
