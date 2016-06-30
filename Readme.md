
  `replay-logs(1)` - is a quick utility that reads cloudfront logs from stdin and replays them
  to a configurable destination, with concurrency and max reqs per second.

  This is useful if you want to stress test with real-data just before deployment.

## Example

  First read a sample of logs from S3 and save them locally:

  ```bash
  $ aws s3 ls s3://<bucket>/cloudfront/<distribution>.<date> \
    | awk '{ printf("s3://<bucket>/cloudfront/%s", $4) }' \
    | xargs -I % aws s3 cp % \
    | gzcat > logs.txt
  ```

  Next run the command, supplying an the address:

  ```bash
  $ replay-logs --addr http://stress.foo.baz < logs.txt
  ```

## Usage

```bash

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


```
