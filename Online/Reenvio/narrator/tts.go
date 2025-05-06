package narrator

import (
	"log"
	"net/http"
	"net/url"
)

func Say(s string) {
	baseURL := "http://tts.docker:3000/"

	params := url.Values{}
	params.Add("text", s)

	finalURL := baseURL + "?" + params.Encode()
	log.Println("finalURL:", finalURL)

	resp, err := http.Get(finalURL)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("HTTP error status:", resp.Status)
		return
	}

	log.Println("Request successfully sent")
}
