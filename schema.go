package speedtest

import (
	"net"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) (err error) {
	v, err := time.Parse(`"2006-01-02T15:04:05Z"`, string(b))
	if err != nil {
		return err
	}
	t.Time = v
	return
}

type Message struct {
	Type      string    `json:"type"`
	Timestamp Timestamp `json:"timestamp"`
}

type Iface struct {
	InternalIP net.IP `json:"internalIp"`
	Name       string `json:"name"`
	MACAddr    string `json:"macAddr"`
	IsVPN      bool   `json:"isVpn"`
	ExternalIP net.IP `json:"externalIp"`
}

type Server struct {
	Id       int    `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Country  string `json:"country"`
	IP       net.IP `json:"ip"`
}

type Info struct {
	ISP    string `json:"isp"`
	Iface  Iface  `json:"interface"`
	Server Server `json:"server"`
}

type PingProgress struct {
	Progress float64 `json:"progress"`
	Jitter   float64 `json:"jitter"`
	Latency  float64 `json:"latency"`
}

type TransmissionProgress struct {
	Bandwidth int64    `json:"bandwidth"`
	Bytes     int64    `json:"bytes"`
	Elapsed   int64    `json:"elapsed"`
	Progress  float64  `json:"progress"`
	Latency   *Latency `json:"latency,omitempty"`
}

type Latency struct {
	Iqm float64 `json:"iqm"`
}

type ResultPing struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
}

type ResultTransmission struct {
	Bandwidth int64   `json:"bandwidth"`
	Bytes     int64   `json:"bytes"`
	Elapsed   int64   `json:"elapsed"`
	Latency   Latency `json:"latency"`
}

type ResultLatency struct {
	Iqm    float64 `json:"iqm"`
	Jitter float64 `json:"jitter"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
}

type Result struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Persisted bool   `json:"persisted"`
}

type StartMessage struct {
	Message
	Info
}

type PingMessage struct {
	Message
	Ping PingProgress `json:"ping"`
}

type DownloadMessage struct {
	Message
	Download TransmissionProgress `json:"download"`
}

type UploadMessage struct {
	Message
	Upload TransmissionProgress `json:"upload"`
}

type ResultMessage struct {
	Message
	Info
	Ping     ResultPing         `json:"ping"`
	Download ResultTransmission `json:"download"`
	Upload   ResultTransmission `json:"upload"`
	Result   Result             `json:"result"`
}
