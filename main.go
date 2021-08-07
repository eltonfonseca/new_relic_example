package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	relic "new_relic_example/new_relic"
	"os"
	"path/filepath"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Host struct {
	Name string
}

func main() {
	appName := "App Example"
	license := "LICENSE KEY"

	newRelicClient := relic.New(appName, license)
	txn := newRelicClient.StartTransaction("transaction")

	hosts := readFile(txn)
	testHosts(hosts, txn)

	txn.End()
}

func testHosts(hosts []Host, txn *newrelic.Transaction) {
	for _, host := range hosts {
		segment := txn.StartSegment("testHosts")

		if host.Name != "" {
			response, err := http.Get(host.Name)

			if err != nil {
				txn.NoticeError(newrelic.Error{
					Message: "error on make request",
					Class:   "TestHosts",
					Attributes: map[string]interface{}{
						"HOSTNAME": host.Name,
					},
				})

				fmt.Println("error on make request")
			}

			segment.AddAttribute("HOSTNAME", host.Name)
			segment.AddAttribute("STATUS", response.StatusCode)

			if response.StatusCode == 200 {
				fmt.Printf("Host: %s está online\n", host.Name)
			} else {
				fmt.Printf("Host: %s está offline\n", host.Name)
			}
		}

		segment.End()
	}
}

func readFile(txn *newrelic.Transaction) []Host {
	var hosts []Host

	path, _ := filepath.Abs("hosts.txt")
	file, err := os.Open(path)

	if err != nil {
		txn.NoticeError(newrelic.Error{
			Message: "error on open file",
			Class:   "ReadFile",
		})

		fmt.Println("error on open file")
		return nil
	}

	reader := bufio.NewReader(file)

	for {
		row, err := reader.ReadString('\n')
		row = strings.TrimSpace(row)

		host := Host{Name: row}
		hosts = append(hosts, host)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return hosts
}
