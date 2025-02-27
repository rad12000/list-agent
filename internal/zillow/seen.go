package zillow

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	writeIgnoreCh = make(chan string)
	writeMatchCh  = make(chan string)
	readSeenCh    = make(chan pair[chan bool, string])
	persistence   *os.File
)

func startListening(file string) {
	persistence = Must(os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644))
	go listenSeen()
}

func listenSeen() {
	m := make(map[string]bool)
	reader := bufio.NewReader(persistence)
	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		parts := strings.Split(string(line), " ")
		m[parts[0]] = parts[1] == "true"
	}

	for {
		var err error
		select {
		case v := <-writeIgnoreCh:
			m[v] = false
			_, err = persistence.WriteString(fmt.Sprintf("%s false\n", v))
		case v := <-writeMatchCh:
			m[v] = true
			_, err = persistence.WriteString(fmt.Sprintf("%s true\n", v))
			cobra.CheckErr(exec.Command("open", "-a", "Google Chrome", v).Run())
		case tuple := <-readSeenCh:
			_, ok := m[tuple.Second]
			tuple.First <- ok
		}

		if err != nil {
			panic(err)
		}
	}
}

func markAsSeen(uri string, matched bool) {
	if matched {
		writeMatchCh <- uri
	} else {
		writeIgnoreCh <- uri
	}
}

func hasBeenSeen(uri string) bool {
	ch := make(chan bool)
	readSeenCh <- pair[chan bool, string]{ch, uri}
	return <-ch
}
