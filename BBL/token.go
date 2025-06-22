package bbl

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func GetTokenFromUrl(url, userName, password string) (string, error) {
	credentials := map[string]string{
		userName: userName,
		password: password,
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(token), nil
}
