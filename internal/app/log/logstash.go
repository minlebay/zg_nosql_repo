package log

import (
	"bytes"
	"net/http"
	"time"
)

type LogstashWriter struct {
	URL string
}

func NewLogstashWriter(c *Config) *LogstashWriter {
	return &LogstashWriter{
		URL: c.Url,
	}
}

func (lw *LogstashWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", lw.URL, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}

	return len(p), nil
}
