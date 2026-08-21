package common

import "sync"

type LogResponse struct {
	Type    string `json:"type"`
	Status  int    `json:"status,omitempty"`
	URI     string `json:"uri,omitempty"`
	Time    int64  `json:"time,omitempty"`
	Message string `json:"message,omitempty"`
}

type logHub struct {
	mu sync.RWMutex

	clients map[chan LogResponse]struct{}
	history []LogResponse
}
