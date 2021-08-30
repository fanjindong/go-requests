package requests

import (
	"context"
	"net/http"
)

type Request struct {
	*http.Request
}

// NewRequest wraps NewRequestWithContext using the background context.
func NewRequest(method, url string) (*Request, error) {
	r, err := http.NewRequestWithContext(context.Background(), method, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", userAgent)
	return &Request{r}, nil
}
