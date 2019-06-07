package stm

import (
	"bytes"
	"compress/gzip"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"time"
)

type SCPAdapter struct {
	sshConfig *ssh.ClientConfig
	client    *scp.Client
	hostname  string
}

func NewSCPAdapter(sshConfig ssh.ClientConfig, hostname string, timeout time.Duration) (*SCPAdapter, error) {
	s := &SCPAdapter{
		sshConfig: &sshConfig,
		hostname:  hostname,
	}
	var client scp.Client
	if timeout == 0 {
		client = scp.NewClient(hostname, &sshConfig)
	} else {
		client = scp.NewClientWithTimeout(hostname, &sshConfig, timeout)
	}
	if err := client.Connect(); err != nil {
		return nil, err
	}
	client.Close()

	return s, nil
}

func (s *SCPAdapter) Write(loc *Location, data []byte) {
	var reader io.Reader = bytes.NewReader(data)

	if GzipPtn.MatchString(loc.Filename()) {
		var writer io.WriteCloser

		reader, writer = io.Pipe()

		gz := gzip.NewWriter(writer)
		_, err := io.Copy(gz, bytes.NewReader(data))
		if err != nil {
			log.Fatal("[F] SCP Upload file Error:", err)
		}
		if err := gz.Close(); err != nil {
			log.Fatal("[F] SCP Upload file Error:", err)
		}
		if err := writer.Close(); err != nil {
			log.Fatal("[F] SCP Upload file Error:", err)
		}
	}
	if err := s.client.Connect(); err != nil {
		log.Fatal("[F] SCP Upload file Error:", err)
	}
	defer s.client.Close()
	if err := s.client.Copy(reader, loc.Filename(), "0666", int64(len(data))); err != nil {
		log.Fatal("[F] SCP Upload file Error:", err)
	}
}

func (s *SCPAdapter) Bytes() [][]byte {
	// TODO
	return nil
}
