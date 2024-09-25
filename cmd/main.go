package main

import (
	"encoding/json"
	"fmt"

	"github.com/lightyen/speedtest-go"
)

func main() {
	t := speedtest.New(speedtest.Options{
		AcceptLicense: true,
		OnStart: func(m *speedtest.StartMessage) {
			fmt.Printf("START %+v\n", m)
		},
		OnDownload: func(m *speedtest.DownloadMessage) {
			fmt.Printf("Download: %.1f Mbps\n", float64(m.Download.Bandwidth<<3)/1e6)
		},
		OnUpload: func(m *speedtest.UploadMessage) {
			fmt.Printf("Upload: %.1f Mbps\n", float64(m.Upload.Bandwidth<<3)/1e6)
		},
		OnResult: func(m *speedtest.ResultMessage) {
			dl, ul := float64(m.Download.Bandwidth<<3)/1e6, float64(m.Upload.Bandwidth<<3)/1e6
			fmt.Printf("DL: %.2f Mbps, UL: %.2f Mbps\n", dl, ul)
			fmt.Printf("Latency: Idle %.0f ms, Download %.0f ms, Upload %.0f ms\n", m.Ping.Latency, m.Download.Latency.Iqm, m.Upload.Latency.Iqm)
			fmt.Printf("%s / %s / %s\n", m.ISP, m.Server.Name, m.Server.Location)
			fmt.Println(m.Result.URL)
			data, _ := json.MarshalIndent(m, "", "  ")
			fmt.Println(string(data))
		},
	})
	err := t.Run()
	if err != nil {
		fmt.Println(err)
	}
}
