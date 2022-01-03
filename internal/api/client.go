package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Client interface {
	SendRequest(req *http.Request) (Response, error)
	Close() error
}

type client struct {
	httpClient *http.Client
	mu         sync.Mutex
	logger     *logrus.Logger
}

func NewClient(logger *logrus.Logger) Client {
	hCli := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &client{httpClient: hCli, logger: logger}
}

func (c *client) SendRequest(req *http.Request) (apiResp Response, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return apiResp, errors.New(fmt.Sprintf("failed to send request, error %s", err))
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&apiResp); err != nil {
		c.logger.Errorf("Error decoding response %v", err)
		return apiResp, errors.New("error parsing response")
	}

	if response.StatusCode != http.StatusOK {
		return apiResp, errors.New(fmt.Sprintf("error from reporting app with status code %d and error %s", response.StatusCode, apiResp.Error))
	}

	return apiResp, nil
}

func (c *client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.httpClient = nil
	return nil
}
