package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type logEntry struct {
	ts     int64
	value  string
	labels map[string]string
}

type lokiLogger struct {
	URL        string
	Key        string
	BufferSize int8
	Level      int8
	Labels     map[string]string
	buffer     []logEntry
	mu         sync.Mutex
}

var _ io.WriteCloser = (*lokiLogger)(nil)

func (l *lokiLogger) Close() error {
	return nil
}

func (l *lokiLogger) Write(p []byte) (n int, err error) {
	l.mu.Lock()

	if l.buffer == nil {
		l.buffer = []logEntry{}
	}

	le := logEntry{
		ts:    time.Now().UnixNano(),
		value: string(p),
	}
	l.buffer = append(l.buffer, le)

	if len(l.buffer) >= int(l.BufferSize) {
		return l.push()
	}

	l.mu.Unlock()
	return 0, nil
}

type line []string

type stream struct {
	Stream map[string]string `json:"stream,omitempty"`
	Values []line            `json:"values,omitempty"`
}

type payload struct {
	Streams []stream `json:"streams,omitempty"`
}

func (l *lokiLogger) push() (n int, err error) {
	lines := []line{}
	for _, v := range l.buffer {
		line := line{}
		line = append(line, strconv.FormatInt(v.ts, 10))
		line = append(line, v.value)
		lines = append(lines, line)
	}

	l.buffer = []logEntry{}
	l.mu.Unlock()

	payload := payload{
		Streams: []stream{
			{
				Stream: l.Labels,
				Values: lines,
			},
		},
	}

	out, err := json.Marshal(&payload)
	if err != nil {
		return 0, err
	}

	bodyReader := bytes.NewReader([]byte(out))

	requestURL := l.URL
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return 0, nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+l.Key)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		fmt.Printf("client: error making http request: %d\n", res.StatusCode)
		return 0, errors.New(fmt.Sprint(res.StatusCode))
	}

	return 0, nil
}
