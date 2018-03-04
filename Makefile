GOOS ?= linux

default: scp

transfer_plex: main.go
	CGO_ENABLED=0 GOOS=$(GOOS) go build -a -installsuffix cgo -o $@ $<

scp: transfer_plex
	scp $< plex@minty:/var/lib/plexmediaserver

test:
	go test . -v -split "roy_rogers/" -src.server roy@rogers.biz

.PHONY: test
