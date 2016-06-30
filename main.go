package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/segmentio/replay-logs/internal/worker"
	"github.com/tj/docopt"
)

const version = ""
const usage = `
  Usage: replay-logs
    [--addr host]
    [--concurrency n]
    [--headers v...]
    [--rate n]

  Example:

    replay-logs --concurrency 10 < logs.txt
    replay-logs -H Accept-Encoding:gzip < logs.txt

  Options:
    --concurrency n    request concurrency [default: 10]
    --addr addr        addr to use, the log path will be appended [default: http://localhost:80]
    --rate n           max requests per second [default: 100]
    -H, --headers v    headers to add to the request [default: ]
    -h, --help         show help information
    -v, --version      show version information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	headers, err := parseHeaders(args["--headers"].([]string))
	if err != nil {
		log.Fatal(err)
	}

	concurrency, err := strconv.Atoi(args["--concurrency"].(string))
	if err != nil {
		log.Fatal(err)
	}

	rate, err := strconv.Atoi(args["--rate"].(string))
	if err != nil {
		log.Fatal(err)
	}

	worker := worker.New(worker.Config{
		Addr:        args["--addr"].(string),
		Concurrency: concurrency,
		Rate:        rate,
		Input:       os.Stdin,
		Headers:     headers,
	})

	err = worker.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func parseHeaders(headers []string) (map[string]string, error) {
	ret := make(map[string]string)

	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", header)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		ret[key] = val
	}

	return ret, nil
}
