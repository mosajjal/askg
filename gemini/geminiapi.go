package gemini

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=YOUR_API_KEY
const geminiAPIURL string = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"

// {
// "candidates": [
//
//	{
//	  "content": {
//	    "parts": [
//	      {
//	        "text":
type apiAnswer struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type apiQuetsion struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// GeminiAPI is the main struct for the Gemini AI API
type GeminiAPI struct {
	apiKey string
	logger *zerolog.Logger

	// Timeout in seconds
	TimeoutQuery int
}

// NewWeb creates a new Gemini AI instance. Cookies is all the cookies from a browser session
func NewAPI(l *zerolog.Logger, apiKey string) AI {
	b := &GeminiAPI{
		apiKey:       apiKey,
		logger:       l,
		TimeoutQuery: 60,
	}
	return b
}

// Ask generates a Gemini AI response and returns it to the user
func (b *GeminiAPI) Ask(prompt string) (string, error) {
	// Create a Resty Client
	client := resty.New()

	client.SetLogger(Log{b.logger})
	client.SetDebug(true)

	// set content type to application/json
	client.SetHeader("Content-Type", "application/json")

	// set the GET parameter of key to apikey
	client.SetQueryParam("key", b.apiKey)

	// '{"contents":[{"parts":[{"text":"Write a story about a magic backpack"}]}]}'
	reqStruct := apiQuetsion{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// convert the struct to a JSON string
	reqJSON, err := json.Marshal(reqStruct)
	if err != nil {
		return "", err
	}

	// client.SetBaseURL(geminiAPIURL)
	if b.TimeoutQuery > 0 {
		client.SetTimeout(time.Duration(b.TimeoutQuery) * time.Second)
	}
	resp, err := client.R().
		SetBody(
			reqJSON,
		).
		SetContentLength(true).
		Post(geminiAPIURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("status code is not 200: %d", resp.StatusCode())
	}

	a := apiAnswer{}
	err = json.Unmarshal(resp.Body(), &a)
	if err != nil {
		return "", err
	}

	if len(a.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}
	if len(a.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in response")
	}
	// get the main answer
	answer := a.Candidates[0].Content.Parts[0].Text
	return answer, nil
}
