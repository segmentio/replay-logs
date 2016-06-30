
  `replay-logs(1)` - is a quick utility that reads cloudfront logs from stdin and replays them
  to a configurable destination, with concurrency and max reqs per second.
  
<img width="962" alt="screen shot 2016-06-30 at 1 30 29 pm" src="https://cloud.githubusercontent.com/assets/1661587/16485297/e071bca4-3ec6-11e6-989b-abac4238d95a.png">

## Example

  First read a sample of logs from S3 and save them locally:

  ```bash
  $ aws s3 ls s3://<bucket>/cloudfront/<distribution>.<date> \
    | awk '{ printf("s3://<bucket>/cloudfront/%s", $4) }' \
    | xargs -I % aws s3 cp % - \
    | gzcat > logs.txt
  ```

  Next run the command, supplying the destination address:

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

## License

Released under the MIT License

(The MIT License)

Copyright (c) 2016 Segment friends@segment.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
