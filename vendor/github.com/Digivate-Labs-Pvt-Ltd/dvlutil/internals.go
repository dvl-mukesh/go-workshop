package dvlutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func sendPostRequestWithBody(reqBody map[string]interface{}, queryURL string, headers map[string]string, options ...HttpOption) (*ApiResponse, error) {
	var reqBytes bytes.Buffer
	err := json.NewEncoder(&reqBytes).Encode(reqBody)

	if err != nil {
		return nil, err
	}

	return sendRequest(queryURL, reqBytes, headers, options...)

}

func sendRequest(queryURL string, reqBytes bytes.Buffer, headers map[string]string, options ...HttpOption) (*ApiResponse, error) {

	req, err := http.NewRequest(http.MethodPost, queryURL, &reqBytes)

	if err != nil {
		return nil, err
	}

	addHeaders(req, headers)

	httpConfig := DefaultHttpConfig()

	for _, option := range options {
		option(httpConfig)
	}

	resp, err := httpConfig.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &ApiResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	if headers == nil {
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
}
