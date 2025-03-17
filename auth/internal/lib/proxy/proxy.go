package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewProxyRequest (targetURL string, requestBody any, method string) (*http.Request, error) {
	if method == "" {
		return nil, fmt.Errorf("empty method")
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	var body *bytes.Buffer
	if len(reqBody) > 0 {
		body = bytes.NewBuffer(reqBody)
	} else if len(reqBody) == 0 {
		body = bytes.NewBuffer([]byte{})
	}

	proxy, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return nil, err
	}

	return proxy, nil
}