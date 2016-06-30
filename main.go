package main

import (
	"log"
	"os"
	"strconv"

	"github.com/segmentio/replay-logs/internal/worker"
	"github.com/tj/docopt"
)

const version = ""
const usage = `
  Usage: replay-logs
    [--addr host]
    [--concurrency n]
    [--rate n]

  Options:
    --concurrency n    request concurrency [default: 10]
    --addr addr        addr to use, the log path will be appended [default: http://localhost:80]
    --rate n           max requests per second [default: 100]
    -h, --help         show help information
    -v, --version      show version information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
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
	})

	err = worker.Run()
	if err != nil {
		log.Fatal(err)
	}
}
