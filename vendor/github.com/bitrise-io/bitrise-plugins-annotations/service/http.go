package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	envServiceURL  = "BITRISEIO_BUILD_ANNOTATIONS_SERVICE_URL"
	envAccessToken = "BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN"
	envBuildSlug   = "BITRISE_BUILD_SLUG"

	defaultServiceURL = "https://build-annotations.services.bitrise.io"
)

func do(method, route string, body any) ([]byte, error) {
	token := os.Getenv(envAccessToken)
	if token == "" {
		return nil, errors.New("access token not found")
	}

	req, err := newRequest(method, route, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	return io.ReadAll(resp.Body)
}

func newRequest(method, route string, body any) (*http.Request, error) {
	build := os.Getenv(envBuildSlug)
	if build == "" {
		return nil, errors.New("build not found")
	}

	serviceURL := os.Getenv(envServiceURL)
	if serviceURL == "" {
		serviceURL = defaultServiceURL
	}

	url := fmt.Sprintf("%s/v1/builds/%s/annotations%s", serviceURL, build, route)

	if body == nil {
		return http.NewRequest(method, url, nil)
	}

	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
