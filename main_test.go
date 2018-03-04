package main

import (
	"context"
	"fmt"
	"testing"
)

func TestURLUnescaping(t *testing.T) {
	r := &request{
		Destination: "tv",
		Link:        "sftp://example.biz/mnt/mpathm/roy_rogers/files/Blade%20Runner%202049%201080p%20WEB-DL%20H264%20AC3-EVO",
	}

	fmt.Println(r.path())
	fmt.Println(r.sftpCmd(context.Background()).Args)

	t.Fail()
}
