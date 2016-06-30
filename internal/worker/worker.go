package worker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/go-stats"
	"github.com/tj/go-sync/semaphore"
)

// Log is a cloudfront Log.
type log struct {
	method string
	path   string
	status int
}

// Config is the configuration.
type Config struct {
	Input       io.Reader
	Headers     map[string]string
	Addr        string
	Rate        int
	Concurrency int
}

// Worker is a request worker.
type Worker struct {
	scanner *bufio.Scanner
	headers map[string]string
	addr    string
	logc    chan log
	rate    *time.Ticker
	sema    semaphore.Semaphore
	stats   *stats.Stats
}

// New returns a new *Worker.
func New(c Config) *Worker {
	w := &Worker{
		scanner: bufio.NewScanner(c.Input),
		logc:    make(chan log),
		sema:    make(semaphore.Semaphore, c.Concurrency),
		rate:    time.NewTicker(time.Second / time.Duration(c.Rate)),
		stats:   stats.New(),
		headers: c.Headers,
		addr:    c.Addr,
	}

	go w.request()
	return w
}

// Run runs the worker.
func (w *Worker) Run() error {
	go w.stats.TickEvery(time.Second)

	for w.scanner.Scan() {
		line := w.scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		log, err := parse(line)
		if err != nil {
			return err
		}

		w.logc <- log
	}

	w.rate.Stop()
	w.sema.Wait()
	close(w.logc)
	w.stats.Stop()
	return w.scanner.Err()
}

// Request runs requests for each log.
func (w *Worker) request() {
	for log := range w.logc {
		<-w.rate.C

		w.sema.Run(func() {
			url := w.addr + log.path

			req, err := http.NewRequest(log.method, url, nil)
			if err != nil {
				w.errorf("unable to create request: %s", err)
				return
			}

			for key, value := range w.headers {
				req.Header.Set(key, value)
			}

			resp, err := http.DefaultClient.Do(req)

			if resp != nil {
				_, _ = io.Copy(ioutil.Discard, resp.Body)
				_ = resp.Body.Close()
			}

			if err != nil {
				w.errorf("request error: %s", err)
				return
			}

			w.stats.Incr("requests")
			w.stats.Incr(fmt.Sprintf("responses.%d", resp.StatusCode))
		})
	}
}

// Log an error.
func (w *Worker) errorf(s string, args ...interface{}) {
	l := fmt.Sprintf(s, args...)
	fmt.Fprintf(os.Stderr, "replay-log: %s\n", l)
}

// Parse parses the line parts into a log.
func parse(line string) (log log, err error) {
	keys := strings.Split(line, "\t")

	if len(keys) != 23 {
		return log, fmt.Errorf("invalid log line: %s", line)
	}

	log.method = keys[5]
	log.path = keys[7]
	log.status, err = strconv.Atoi(keys[8])
	if err != nil {
		return log, fmt.Errorf("invalid status code %s (%s)", err, keys[8])
	}

	return log, nil
}
