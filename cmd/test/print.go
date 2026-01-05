package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	RESET = "\x1b[0m"
	RED   = "\x1b[31m"
)

func print(text string) bool {
	if flags.Format != "text" {
		return false
	}
	fmt.Println(text)
	return true
}

func printf(format string, a ...any) bool {
	if flags.Format != "text" {
		return false
	}
	fmt.Println(fmt.Sprintf(format, a...))
	return true
}

func printFormattedResults() bool {
	switch flags.Format {
	case "text":
		return false
	case "json":
		return printJson(true)
	case "jsonl":
		return printJson(false)
	case "influx":
		return printInflux()
	case "tsv":
		return printTsv()
	}
	stderrRed("Unsupported format: %s", flags.Format)
	return false
}

func stderr(format string, a ...any) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, a...))
}

func isTtyStdout() bool {
	o, err := os.Stdout.Stat()
	return err == nil && (o.Mode()&os.ModeCharDevice) == os.ModeCharDevice
}

func isTtyStderr() bool {
	o, err := os.Stderr.Stat()
	return err == nil && (o.Mode()&os.ModeCharDevice) == os.ModeCharDevice
}

func stderrRed(format string, a ...any) {
	if isTtyStderr() {
		// Colorful output to the terminal
		fmt.Fprintln(os.Stderr, RED+fmt.Sprintf(format, a...)+RESET)
		return
	}

	// Output to pipe, no colors
	stderr(format, a...)
}

func progressln(format string, a ...any) {
	if flags.NoProgress {
		return
	}
	stderr(format, a...)
}

func progress(format string, a ...any) {
	if flags.NoProgress {
		return
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

func printJson(indent bool) bool {
	if indent {
		result, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return false
		}
		fmt.Println(string(result))
		return true
	}

	result, err := json.Marshal(result)
	if err != nil {
		return false
	}
	fmt.Println(string(result))
	return true
}

func printInflux() bool {
	// InfluxDB Line Protocol uses only `\n` for line separation
	// So we’re not using Println-like functions here
	fmt.Print("#####\n")
	fmt.Print("# InfluxDB Line Protocol\n")
	fmt.Print("#\n")
	fmt.Printf("# @time %s\n", result.Time.Format(time.RFC3339))
	fmt.Print("# @metrics 11\n")
	fmt.Print("#####\n")
	fmt.Print("\n")

	// Tags are safe values, so no escaping needed
	tags := fmt.Sprintf(
		"asn=%d,ip=%s,id=%s",
		result.Connection.Asn,
		result.Connection.Ip,
		result.Id,
	)
	fmt.Printf("#\n")
	fmt.Printf("# UNIT speed_download bits per second\n")
	fmt.Printf("speed_download,%s        value=%du\n", tags, result.Speed.Download)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT speed_upload bits per second\n")
	fmt.Printf("speed_upload,%s          value=%du\n", tags, result.Speed.Upload)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT latency_unloaded milliseconds\n")
	fmt.Printf("latency_unloaded,%s      value=%f\n", tags, result.Latency.Unloaded.Median)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT jitter_unloaded milliseconds\n")
	fmt.Printf("jitter_unloaded,%s       value=%f\n", tags, result.Jitter.Unloaded)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT latency_downloaded milliseconds\n")
	fmt.Printf("latency_downloaded,%s    value=%f\n", tags, result.Latency.Downloaded.Median)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT jitter_downloaded milliseconds\n")
	fmt.Printf("jitter_downloaded,%s     value=%f\n", tags, result.Jitter.Downloaded)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT latency_uploaded milliseconds\n")
	fmt.Printf("latency_uploaded,%s      value=%f\n", tags, result.Latency.Uploaded.Median)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT jitter_uploaded milliseconds\n")
	fmt.Printf("jitter_uploaded,%s       value=%f\n", tags, result.Jitter.Uploaded)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT bytes_downloaded bytes\n")
	fmt.Printf("bytes_downloaded,%s     value=%d\n", tags, result.Traffic.TotalDownloadedBytes)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT bytes_uploaded bytes\n")
	fmt.Printf("bytes_uploaded,%s     value=%d\n", tags, result.Traffic.TotalUploadedBytes)

	fmt.Printf("#\n")
	fmt.Printf("# UNIT test_duration seconds\n")
	fmt.Printf("test_duration,%s       value=%f\n", tags, result.Duration)

	fmt.Print("\n")

	stderr("# NOTE: InfluxDB Line Protocol not yet fully supported")
	return true
}

func printTsv() bool {
	columns := []string{
		"Test Time",
		"Download Speed",
		"Upload Speed",
		"Unloaded Latency",
		"Unloaded Jitter",
		"Downloaded Latency",
		"Downloaded Jitter",
		"Downloaded Bytes",
		"Uploaded Latency",
		"Uploaded Jitter",
		"Uploaded Bytes",
		"Test Duration",
		"ASN",
		"IP",
		"Test ID",
	}
	data := []string{
		result.Time.Format(time.RFC3339),
		fmt.Sprintf("%d", result.Speed.Download),
		fmt.Sprintf("%d", result.Speed.Upload),
		fmt.Sprintf("%f", result.Latency.Unloaded.Median),
		fmt.Sprintf("%f", result.Jitter.Unloaded),
		fmt.Sprintf("%f", result.Latency.Downloaded.Median),
		fmt.Sprintf("%f", result.Jitter.Downloaded),
		fmt.Sprintf("%d", result.Traffic.TotalDownloadedBytes),
		fmt.Sprintf("%f", result.Latency.Uploaded.Median),
		fmt.Sprintf("%f", result.Jitter.Uploaded),
		fmt.Sprintf("%d", result.Traffic.TotalUploadedBytes),
		fmt.Sprintf("%f", result.Duration),
		fmt.Sprintf("%d", result.Connection.Asn),
		result.Connection.Ip,
		result.Id,
	}
	fmt.Println(strings.Join(columns, "\t"))
	fmt.Println(strings.Join(data, "\t"))
	return true
}
