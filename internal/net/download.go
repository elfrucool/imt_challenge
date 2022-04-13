package net

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Open the connection
func StartDownload(urlString string) (io.ReadCloser, error) {
	client := http.Client{}

	response, err := client.Get(urlString)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("Error downloading file: %v", response.Status))
	}

	return response.Body, nil
}
