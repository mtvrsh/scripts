package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// https://github.com/ewhal/time-distance-golang
// rewrite into dateutil? (include timediff)
func main() {
	var r io.Reader = os.Stdin
	if len(os.Args) > 1 {
		var buff bytes.Buffer
		for _, arg := range os.Args[1:] {
			fmt.Fprintln(&buff, arg)
		}
		r = &buff
	}

	for sc := bufio.NewScanner(r); sc.Scan(); {
		dateStr := strings.TrimSpace(sc.Text())
		// p, err := dateparse.ParseStrict(dateStr)
		p, err := dateparse.ParseAny(dateStr)
		cerr(err)

		fmt.Printf("%d\t%d\t%d\t%s\n", p.UnixNano(), p.UnixMilli(), p.Unix(), p.Format(time.RFC3339))
	}
}

func cerr(err error) {
	if err != nil {
		panic(err)
	}
}
