// Package haste is a hastebin client.
package haste

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Response contains a haste response, which is just a key.
type Response struct {
	Key string `json:"key"`
	Message string `json:"message"`
}

// GetLink returns a full URL to a hastebin key.
// This requires the Haste instance to get the host.
func (resp *Response) GetLink(haste *Haste) string {
	return haste.Host + "/" + resp.Key
}

// Haste is a Hastebin client.
type Haste struct {
	Host string
}

// NewHaste creates a new Haste instance with a specified URL basepoint.
func NewHaste(host string) *Haste {
	return &Haste{
		Host: host,
	}
}

// Fetch gets an already existing item from Hastebin
func (haste *Haste) Fetch(key string) (string, error) {
	resp, err := http.Get(haste.Host + "/raw/" + key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// UploadString uploads a string to Hastebin.
func (haste *Haste) UploadString(data string) (*Response, error) {
	return haste.UploadBuffer(bytes.NewBuffer([]byte(data)))
}

// UploadBytes uploads bytes to Hastebin.
func (haste *Haste) UploadBytes(data []byte) (*Response, error) {
	return haste.UploadBuffer(bytes.NewBuffer(data))
}

// UploadBuffer uploads a buffer to Hastebin.
func (haste *Haste) UploadBuffer(data *bytes.Buffer) (*Response, error) {
	req, err := http.NewRequest("POST", haste.Host+"/documents", data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp Response
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}
	if apiResp.Message == "Document exceeds maximum length." {
		return nil, errors.New("file too large")
	}

	return &apiResp, nil
}
