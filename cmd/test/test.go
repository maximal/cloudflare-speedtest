package test

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/maximal/cloudflare-speedtest/internal/exit"
	"github.com/maximal/cloudflare-speedtest/internal/format"
	"github.com/montanaflynn/stats"
	"github.com/oklog/ulid/v2"
)

var (
	result      Result
	timingRegex = regexp.MustCompile(`cfRequestDuration;dur=([\d\.]+)`)
)

func run() exit.Status {
	result.Id = strings.ToLower(ulid.Make().String())
	result.Time = time.Now()
	result.Info = SHORT_DESCRIPTION

	print(result.Info)
	print("")

	if flags.Insecure {
		flags.BaseUrl = "http://" + SPEED_BASE_URL
	} else {
		flags.BaseUrl = "https://" + SPEED_BASE_URL
	}

	client := &http.Client{
		Timeout: 40 * time.Second,
		// Transport: http.RoundTripper{},
	}

	progressln("Getting connection info...")

	result.Latencies.Unloaded = make([]float64, 0, flags.LatencyCount)
	result.Latencies.Downloaded = make([]float64, 0)
	result.Latencies.Uploaded = make([]float64, 0)

	result.Jitters.Unloaded = make([]float64, 0, flags.LatencyCount-1)
	result.Jitters.Downloaded = make([]float64, 0)
	result.Jitters.Uploaded = make([]float64, 0)

	result.Speeds.Download = make([]uint64, 0)
	result.Speeds.Upload = make([]uint64, 0)

	metrics, response, err := getResponseStats(client, false, 0)
	if err != nil {
		log.Fatal(err)
	}

	result.Latencies.Unloaded = append(result.Latencies.Unloaded, metrics.LatencyMs)
	result.Tests.Latency++
	lastLatency := metrics.LatencyMs

	result.Traffic.TotalDownloadedBytes += metrics.ResponseSize
	result.Traffic.TotalUploadedBytes += metrics.RequestSize

	result.Connection.Asn, _ = strconv.ParseUint(response.Header.Get("Cf-Meta-Asn"), 10, 64)
	result.Connection.AsnLink = ASN_BASE_URL + response.Header.Get("Cf-Meta-Asn")
	result.Connection.Ip = response.Header.Get("Cf-Meta-Ip")
	result.Connection.County = response.Header.Get("Cf-Meta-Country")
	result.Connection.City = response.Header.Get("Cf-Meta-City")
	result.Connection.Latitude, _ = strconv.ParseFloat(
		response.Header.Get("Cf-Meta-Latitude"),
		64,
	)
	result.Connection.Longitude, _ = strconv.ParseFloat(
		response.Header.Get("Cf-Meta-Longitude"),
		64,
	)
	result.Connection.Timezone = response.Header.Get("Cf-Meta-Timezone")

	printf("ASN:         %d    %s", result.Connection.Asn, result.Connection.AsnLink)
	printf("IP:          %s", result.Connection.Ip)
	printf(
		"Location:    %s / %s",
		result.Connection.County,
		result.Connection.City,
	)
	printf("Coordinates: %f, %f", result.Connection.Latitude, result.Connection.Longitude)
	printf("Timezone:    %s", result.Connection.Timezone)
	print("")

	progressln("Running %d latency tests...", flags.LatencyCount)

	for i := range flags.LatencyCount - 1 {
		progress(" %d ", i+2)

		metrics, _, err := getResponseStats(client, false, 0)
		if err != nil {
			log.Fatal(err)
		}

		result.Traffic.TotalDownloadedBytes += metrics.ResponseSize
		result.Traffic.TotalUploadedBytes += metrics.RequestSize

		result.Latencies.Unloaded = append(result.Latencies.Unloaded, metrics.LatencyMs)
		result.Jitters.Unloaded = append(
			result.Jitters.Unloaded,
			math.Abs(metrics.LatencyMs-lastLatency),
		)
		result.Tests.Latency++
		lastLatency = metrics.LatencyMs
	}
	progressln("")

	// var lastLatencyUploaded float64 = 0
	// var lastLatencyDownloaded float64 = 0

	for _, step := range testSteps {
		if step.Upload {
			progressln(
				"Uploading %s × %d times...",
				format.BytesSi(uint64(step.Bytes)),
				step.Count,
			)
		} else {
			progressln(
				"Downloading %s × %d times...",
				format.BytesSi(uint64(step.Bytes)),
				step.Count,
			)
		}

		softLimit := false

		for i := range step.Count {
			if flags.SoftLimit > 0 {
				elapsed := time.Since(result.Time)
				if elapsed > flags.SoftLimit {
					progress(" soft time limit reached ")
					softLimit = true
					break
				}
			}

			progress(" %d ", i+1)

			metrics, _, err := getResponseStats(client, step.Upload, uint64(step.Bytes))
			if err != nil {
				log.Fatal(err)
			}

			result.Traffic.TotalDownloadedBytes += metrics.ResponseSize
			result.Traffic.TotalUploadedBytes += metrics.RequestSize

			if step.Upload {
				result.Speeds.Upload = append(result.Speeds.Upload, metrics.UploadSpeed)
				result.Tests.Upload++

				//if lastLatencyUploaded > 0 {
				//	result.Latencies.Uploaded = append(result.Latencies.Uploaded, metrics.LatencyMs)
				//	result.Jitters.Uploaded = append(
				//		result.Jitters.Uploaded,
				//		math.Abs(metrics.LatencyMs-lastLatencyUploaded),
				//	)
				//}
				//lastLatencyUploaded = metrics.LatencyMs
				//result.Tests.Latency++
			} else {
				result.Speeds.Download = append(result.Speeds.Download, metrics.DownloadSpeed)
				result.Tests.Download++

				//if lastLatencyDownloaded > 0 {
				//	result.Latencies.Downloaded = append(result.Latencies.Downloaded, metrics.LatencyMs)
				//	result.Jitters.Downloaded = append(
				//		result.Jitters.Downloaded,
				//		math.Abs(metrics.LatencyMs-lastLatencyDownloaded),
				//	)
				//}
				//lastLatencyDownloaded = metrics.LatencyMs
				//result.Tests.Latency++
			}
		}
		progressln("")
		if softLimit {
			break
		}
	}

	progressln("Tests done.")
	if isTtyStdout() == isTtyStderr() {
		// If stdout and stderr are both printed to the terminal or a file/pipe,
		// separate progress and results with the blank line
		progressln("")
	}

	// Latencies
	result.Latency.Unloaded = getDataStats(result.Latencies.Unloaded)
	result.Latency.Downloaded = getDataStats(result.Latencies.Downloaded)
	result.Latency.Uploaded = getDataStats(result.Latencies.Uploaded)

	// Jitters
	result.Jitter.Unloaded = getDataStats(result.Jitters.Unloaded).Mean
	result.Jitter.Downloaded = getDataStats(result.Jitters.Downloaded).Mean
	result.Jitter.Uploaded = getDataStats(result.Jitters.Uploaded).Mean

	// Download Speed, 90th percentile
	arr := stats.LoadRawData(result.Speeds.Download)
	if percentile, err := arr.Percentile(90); err == nil {
		result.Speed.Download = uint64(percentile)
	}

	// Upload Speed, 90th percentile
	arr = stats.LoadRawData(result.Speeds.Upload)
	if percentile, err := arr.Percentile(90); err == nil {
		result.Speed.Upload = uint64(percentile)
	}

	printf(
		"Download:    %s    %s/s    %s/s",
		format.BitsPerSecondSi(result.Speed.Download),
		format.BytesSi(result.Speed.Download/8),
		format.BytesIec(result.Speed.Download/8),
	)
	printf(
		"Upload:      %s    %s/s    %s/s",
		format.BitsPerSecondSi(result.Speed.Upload),
		format.BytesSi(result.Speed.Upload/8),
		format.BytesIec(result.Speed.Upload/8),
	)
	printf(
		"Latency:     %.3f ms median    min: %.3f    25p: %.3f    75p: %.3f    max: %.3f",
		result.Latency.Unloaded.Median,
		result.Latency.Unloaded.Min,
		result.Latency.Unloaded.Perc25,
		result.Latency.Unloaded.Perc75,
		result.Latency.Unloaded.Max,
	)
	//printf(
	//	"             During Download:  %.3f ms median    min: %.3f    25p: %.3f    75p: %.3f    max: %.3f",
	//	result.Latency.Downloaded.Median,
	//	result.Latency.Downloaded.Min,
	//	result.Latency.Downloaded.Perc25,
	//	result.Latency.Downloaded.Perc75,
	//	result.Latency.Downloaded.Max,
	//)
	//printf(
	//	"             During Upload:    %.3f ms median    min: %.3f    25p: %.3f    75p: %.3f    max: %.3f",
	//	result.Latency.Uploaded.Median,
	//	result.Latency.Uploaded.Min,
	//	result.Latency.Uploaded.Perc25,
	//	result.Latency.Uploaded.Perc75,
	//	result.Latency.Uploaded.Max,
	//)
	//printf(
	//	"Jitter:      %.3f ms    %.3f ms during download    %.3f ms during upload",
	//	result.Jitter.Unloaded,
	//	result.Jitter.Downloaded,
	//	result.Jitter.Uploaded,
	//)
	printf("Jitter:      %.3f ms", result.Jitter.Unloaded)
	printf(
		"Traffic:     %s / %s downloaded    %s / %s uploaded",
		format.BytesSi(result.Traffic.TotalDownloadedBytes),
		format.BytesIec(result.Traffic.TotalDownloadedBytes),
		format.BytesSi(result.Traffic.TotalUploadedBytes),
		format.BytesIec(result.Traffic.TotalUploadedBytes),
	)
	printf(
		"Tests:       %d latency    %d download    %d upload",
		result.Tests.Latency,
		result.Tests.Download,
		result.Tests.Upload,
	)

	result.Duration = time.Since(result.Time).Seconds()
	printf("Test Time:   %s", result.Time.Format(time.RFC3339))
	printf("Duration:    %f s", result.Duration)
	printf("Test ID:     %s", result.Id)
	print("")
	print(result.Info)

	printFormattedResults()

	return exit.StatusOk
}

func getServerTiming(response *http.Response) time.Duration {
	timing := response.Header.Get("Server-Timing")

	match := timingRegex.FindStringSubmatch(timing)
	if len(match) != 2 {
		return 0
	}
	if val, err := strconv.ParseFloat(match[1], 64); err == nil {
		return time.Duration(val * float64(time.Millisecond))
	}
	return 0
}

func getDataStats(data []float64) Stats {
	if len(data) == 0 {
		return Stats{}
	}
	arr := stats.LoadRawData(data)
	minimum, _ := stats.Min(arr)
	maximum, _ := stats.Max(arr)
	mean, _ := stats.Mean(arr)
	median, _ := stats.Median(arr)
	perc25, err := stats.Percentile(arr, 25)
	if err != nil {
		log.Fatal(err)
	}
	perc75, err := stats.Percentile(arr, 75)
	if err != nil {
		log.Fatal(err)
	}
	return Stats{
		Min:    minimum,
		Perc25: perc25,
		Mean:   mean,
		Median: median,
		Perc75: perc75,
		Max:    maximum,
	}
}

func getResponseHeadersSize(response *http.Response) uint64 {
	// fmt.Println("header")
	// fmt.Println(response.Proto + " " + response.Status)
	headersSize := uint64(len(response.Proto + " " + response.Status + "\r\n"))
	for k, v := range response.Header {
		for _, val := range v {
			// fmt.Printf("%s: %s\n", k, val)
			headersSize += uint64(len(k) + 2 + len(val) + 2)
		}
	}
	return headersSize + 2 // + final \r\n
}

func getRequestTotalSize(request *http.Request, bodySize uint64) uint64 {
	headersSize := uint64(
		len(request.Method + " " + request.URL.String() + " " + request.Proto + "\r\n"),
	)
	for k, v := range request.Header {
		for _, val := range v {
			// fmt.Printf("%s: %s", k, val)
			headersSize += uint64(len(k) + 2 + len(val) + 2)
		}
	}
	if bodySize > 0 {
		return headersSize + 2 + bodySize
	}
	return headersSize
}

func getResponseStats(
	client *http.Client,
	upload bool,
	bodySize uint64,
) (*ResponseStats, *http.Response, error) {
	var method string
	var url string
	var body io.Reader
	var uploadBodySize uint64

	if upload {
		method = "POST"
		url = fmt.Sprintf("%s/__up?measId=%s", flags.BaseUrl, result.Id)
		body = bytes.NewReader(bytes.Repeat([]byte{'0'}, int(bodySize)))
		uploadBodySize = bodySize
	} else {
		method = "GET"
		url = fmt.Sprintf("%s/__down?measId=%s&bytes=%d", flags.BaseUrl, result.Id, bodySize)
		body = nil
		uploadBodySize = 0
	}

	// `connectDone` is now, ConnectDone, or TLSHandshakeDone, whichever comes last
	var connectDone time.Time
	// `uploadDone` is WroteRequest or when request done, whichever comes last
	var uploadDone time.Time
	// `ttfb` is time to first byte
	var ttfb time.Time

	reused := true
	trace := &httptrace.ClientTrace{
		ConnectDone: func(string, string, error) {
			connectDone = time.Now()
			reused = false
			// log.Println("ConnectDone")
			// log.Println(connectDone)
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			connectDone = time.Now()
			reused = false
			// log.Println("TLSHandshakeDone")
			// log.Println(connectDone)
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			// called after the request (including body) has been written
			uploadDone = time.Now()
			// log.Println("uploadDone")
			// log.Println(uploadDone)
		},
		GotFirstResponseByte: func() {
			ttfb = time.Now()
			// log.Println("ttfb")
			// log.Println(ttfb)
		},
	}

	request, _ := http.NewRequestWithContext(
		httptrace.WithClientTrace(context.Background(), trace),
		method,
		url,
		body,
	)
	if upload {
		request.ContentLength = int64(bodySize)
		request.Header.Set("Content-Length", strconv.FormatUint(bodySize, 10))
	}

	connectDone = time.Now()
	// log.Println("connectDone from now")
	// log.Println(connectDone)
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	// log.Println("request done", time.Now())
	uploadDone = time.Now()

	// Downloading body
	bodyWritten, err := io.Copy(io.Discard, response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, nil, err
	}
	_ = response.Body.Close()
	// log.Println("download done", time.Now())

	// Maybe `now - TTFB` is not the best choice
	// downloaded := time.Since(ttfb)
	// Considering `now - connectDone` as better (and lower-numbered) alternative
	downloadDuration := time.Since(connectDone)
	uploadDuration := uploadDone.Sub(connectDone)

	requestSize := getRequestTotalSize(request, uploadBodySize)
	responseSize := getResponseHeadersSize(response) + uint64(bodyWritten)

	serverTiming := getServerTiming(response)
	ttfbDuration := ttfb.Sub(connectDone)
	latency := ttfbDuration - serverTiming
	return &ResponseStats{
		Latency:       latency,
		LatencyMs:     float64(latency) / float64(time.Millisecond),
		Upload:        uploadDuration,
		TTFB:          ttfbDuration,
		Download:      downloadDuration,
		RequestSize:   requestSize,
		ResponseSize:  responseSize,
		UploadSpeed:   uint64(math.Round(float64(requestSize*8) / uploadDuration.Seconds())),
		DownloadSpeed: uint64(math.Round(float64(responseSize*8) / downloadDuration.Seconds())),
		Reused:        reused,
	}, response, nil
}
