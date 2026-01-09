# Cloudflare SpeedTest CLI

This is a command-line version of [Cloudflare SpeedTest](https://speed.cloudflare.com/) utility.

It gets the basic information and metrics about your internet connection.

This tool is not affiliated with Cloudflare, Inc in any way. It only uses Cloudflare testsing servers.


## Build

Get the source code from GitHub:
```shell
git clone https://github.com/maximal/cloudflare-speedtest
cd cloudflare-speedtest
```

Check run using Golang without building:
```shell
go run . --version
# Expected result:
# Cloudflare Speedtest CLI v1.1
```

Build the executable file with Go:
```shell
go build -ldflags '-s -w' -trimpath
```


## Usage

Basic usage with text output (human readable):
```shell
./cloudflare-speedtest
```

Example output:
```plain
Cloudflare Speedtest CLI    © MaximAL    https://github.com/maximal/cloudflare-speedtest

ASN:         6677    https://radar.cloudflare.com/quality/as6677
IP:          66.77.88.99
Provider:    MaximAL ISP
Location:    CZ / Prague
Coordinates: 50.088040, 14.420760
Server:      CZ / Prague / PRG

Download:    536 Mbit/s    67.1 MB/s    63.9 MiB/s
Upload:      111 Mbit/s    13.9 MB/s    13.2 MiB/s
Latency:     8.001 ms median    min: 6.593    25p: 7.333    75p: 8.456    max: 58.900
Jitter:      6.545 ms
Traffic:     594 MB / 567 MiB downloaded    159 MB / 152 MiB uploaded
Tests:       20 latency    33 download    27 upload
Test Time:   2026-01-04T20:54:16+03:00
Duration:    42.559542 s
Test ID:     01ke52agjwv6edvt0xgn7zy4mr
```

### Other Formats
Formatted and indented JSON (for machine use and interchange):
```shell
./cloudflare-speedtest --format=json
# or shorter: ./cloudflare-speedtest -f json
```

[JSON in one line](https://jsonlines.org/) without indentation (safe for appending to JSONL file):
```shell
./cloudflare-speedtest --format=jsonl
```

[InfluxDB line protocol](https://docs.influxdata.com/influxdb3/core/reference/line-protocol/) format output (not yet fully supported):
```shell
./cloudflare-speedtest --format=influx
```

[Tab separated values](https://en.wikipedia.org/wiki/Tab-separated_values) format output:
```shell
./cloudflare-speedtest --format=tsv
```

### Further Help

Run `./cloudflare-speedtest --help` for usage information and instructions:
```plain
Cloudflare Speedtest CLI    © MaximAL    https://github.com/maximal/cloudflare-speedtest

This program tests your internet connection using Cloudflare infrastructure:
https://speed.cloudflare.com/

This tool is not affiliated with Cloudflare, Inc in any way.
It only uses Cloudflare testsing servers.

Contact the author:
* https://github.com/maximal
* https://t.me/maximal
* https://maximals.ru

Usage:
  cloudflare-speedtest [flags]

Flags:
  -f, --format string   output format: text, json, jsonl, influx, tsv (default "text")
  -h, --help            help for cloudflare-speedtest
  -i, --insecure        make HTTP requests instead of HTTPS/TLS ones
  -n, --no-progress     do not print progress information to STDERR
  -v, --version         print version information and exit
```
