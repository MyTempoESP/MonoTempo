package main

import (
	"net/http"
	"net/url"
)

func Say(s string) {
	baseURL := "http://tts.docker:3000/"

	params := url.Values{}
	params.Add("text", s)
	params.Add("speed", "120")

	finalURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(finalURL)

	if err != nil {

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return
	}
}
