package speedtest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrAcceptLicense = errors.New("need to accept license")
	ErrAcceptGDPR    = errors.New("need to accept gdpr")
)

type Options struct {
	AcceptLicense bool
	AcceptGDPR    bool

	OnError    func(m *Message, err error)
	OnStart    func(m *StartMessage)
	OnPing     func(m *PingMessage)
	OnDownload func(m *DownloadMessage)
	OnUpload   func(m *UploadMessage)
	OnResult   func(m *ResultMessage)
}

type SpeedTest struct {
	opts Options
}

func New(options Options) *SpeedTest {
	if options.OnError == nil {
		options.OnError = func(m *Message, err error) {}
	}
	if options.OnStart == nil {
		options.OnStart = func(m *StartMessage) {}
	}
	if options.OnPing == nil {
		options.OnPing = func(m *PingMessage) {}
	}
	if options.OnDownload == nil {
		options.OnDownload = func(m *DownloadMessage) {}
	}
	if options.OnUpload == nil {
		options.OnUpload = func(m *UploadMessage) {}
	}
	if options.OnResult == nil {
		options.OnResult = func(m *ResultMessage) {}
	}
	return &SpeedTest{options}
}

func (s *SpeedTest) Run() error {
	name, err := ensure()
	if err != nil {
		return err
	}

	r, w := io.Pipe()
	var err2 error
	go func() {
		err2 = s.handle(r)
	}()

	args := []string{"-f", "json", "-P", "8", "-p"}

	if s.opts.AcceptLicense {
		args = append(args, "--accept-license")
	}

	if s.opts.AcceptGDPR {
		args = append(args, "--accept-gdpr")
	}

	_, err = execWithOptions(execOptions{
		Stdout: w,
		Stderr: w,
	}, name, args...)

	if err2 == ErrAcceptLicense || err2 == ErrAcceptGDPR {
		return err2
	}

	return err
}

func (s *SpeedTest) handle(r io.Reader) error {
	br := bufio.NewReader(r)

	needAccept := false

	for {
		data, _, err := br.ReadLine()

		if bytes.Index(data, []byte("To accept the message please run speedtest interactively or use the following:")) != -1 {
			needAccept = true
			continue
		}

		if needAccept {
			if bytes.Index(data, []byte("speedtest --accept-license")) != -1 {
				return ErrAcceptLicense
			} else if bytes.Index(data, []byte("speedtest --accept-gdpr")) != -1 {
				return ErrAcceptGDPR
			}

		}

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		var msg *Message
		if err := json.Unmarshal(data, &msg); err != nil {
			s.opts.OnError(msg, err)
			continue
		}

		switch msg.Type {
		case "testStart":
			var p *StartMessage
			if err := json.Unmarshal(data, &p); err != nil {
				s.opts.OnError(msg, err)
				continue
			}
			s.opts.OnStart(p)
		case "ping":
			var p *PingMessage
			if err := json.Unmarshal(data, &p); err != nil {
				s.opts.OnError(msg, err)
				continue
			}
			s.opts.OnPing(p)
		case "download":
			var p *DownloadMessage
			if err := json.Unmarshal(data, &p); err != nil {
				s.opts.OnError(msg, err)
				continue
			}
			s.opts.OnDownload(p)
		case "upload":
			var p *UploadMessage
			if err := json.Unmarshal(data, &p); err != nil {
				s.opts.OnError(msg, err)
				continue
			}
			s.opts.OnUpload(p)
		case "result":
			var p *ResultMessage
			if err := json.Unmarshal(data, &p); err != nil {
				s.opts.OnError(msg, err)
				continue
			}
			s.opts.OnResult(p)
		}
	}
}
