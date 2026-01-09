package test

import "time"

type ResponseStats struct {
	Latency       time.Duration
	LatencyMs     float64
	Upload        time.Duration
	TTFB          time.Duration
	Download      time.Duration
	RequestSize   uint64
	ResponseSize  uint64
	UploadSpeed   uint64
	DownloadSpeed uint64
	Reused        bool
}

type Stats struct {
	Min    float64 `json:"min"`
	Perc25 float64 `json:"perc25"`
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	Perc75 float64 `json:"perc75"`
	Max    float64 `json:"max"`
}

type Result struct {
	Id       string    `json:"id"`
	Time     time.Time `json:"time"`
	Duration float64   `json:"duration"`
	Latency  struct {
		Unloaded   Stats `json:"unloaded"`
		Downloaded Stats `json:"downloaded"`
		Uploaded   Stats `json:"uploaded"`
	} `json:"latency"`
	Jitter struct {
		Unloaded   float64 `json:"unloaded"`
		Downloaded float64 `json:"downloaded"`
		Uploaded   float64 `json:"uploaded"`
	} `json:"jitter"`
	Speed struct {
		Download uint64 `json:"download"`
		Upload   uint64 `json:"upload"`
	} `json:"speed"`
	Latencies struct {
		Unloaded   []float64 `json:"unloaded"`
		Downloaded []float64 `json:"downloaded"`
		Uploaded   []float64 `json:"uploaded"`
	} `json:"latencies"`
	Jitters struct {
		Unloaded   []float64 `json:"unloaded"`
		Downloaded []float64 `json:"downloaded"`
		Uploaded   []float64 `json:"uploaded"`
	} `json:"jitters"`
	Speeds struct {
		Download []uint64 `json:"download"`
		Upload   []uint64 `json:"upload"`
	} `json:"speeds"`
	Traffic struct {
		TotalDownloadedBytes uint64 `json:"total_downloaded_bytes"`
		TotalUploadedBytes   uint64 `json:"total_uploaded_bytes"`
	} `json:"traffic"`
	Tests struct {
		Latency  uint64 `json:"latency"`
		Download uint64 `json:"download"`
		Upload   uint64 `json:"upload"`
	} `json:"tests"`
	Connection struct {
		Asn       uint64  `json:"asn"`
		AsnLink   string  `json:"asn_link"`
		Provider  string  `json:"provider"`
		Ip        string  `json:"ip"`
		Country   string  `json:"country"`
		Region    string  `json:"region"`
		City      string  `json:"city"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"connection"`
	Info string `json:"info"`
}

type ConnectionInfo struct {
	Ip           string `json:"clientIp"`
	Asn          uint64 `json:"asn"`
	Organization string `json:"asOrganization"`
	Country      string `json:"country"`
	City         string `json:"city"`
	Region       string `json:"region"`
	PostalCode   string `json:"postalCode"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Colocation   struct {
		Iata      string  `json:"iata"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
		Country   string  `json:"cca2"`
		Region    string  `json:"region"`
		City      string  `json:"city"`
	} `json:"colo"`
}

func (rs ResponseStats) Debug() {
	printf("Latency: %v", rs.Latency)
	// printf("Latency in ms: %v", rs.LatencyMs)
	printf("Upload: %v", rs.Upload)
	printf("TTFB: %v", rs.TTFB)
	printf("Download: %v", rs.Download)
	printf("Request Size: %v", rs.RequestSize)
	printf("Response Size: %v", rs.ResponseSize)
	printf("Upload Speed: %v", rs.UploadSpeed)
	printf("Download Speed: %v", rs.DownloadSpeed)
	printf("Connection Reused: %v", rs.Reused)
}
