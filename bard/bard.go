package bard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/util/rand"
)

var headers map[string]string = map[string]string{
	"Host":          "bard.google.com",
	"X-Same-Domain": "1",
	"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.4472.114 Safari/537.36",
	"Content-Type":  "application/x-www-form-urlencoded;charset=UTF-8",
	"Origin":        "https://bard.google.com",
	"Referer":       "https://bard.google.com/",
}

const bardURL string = "https://bard.google.com/_/BardChatUi/data/assistant.lamda.BardFrontendService/StreamGenerate"

type bardAnswer struct {
	Content           string   `json:"content"`
	ConversationID    string   `json:"conversationId"`
	ResponseID        string   `json:"responseId"`
	ChoiceID          string   `json:"choiceId"`
	FactualityQueries []string `json:"factualityQueries"`
	TextQuery         string   `json:"textQuery"`
	Choices           []string `json:"choices"`
}

// Bard is the main struct for the Bard AI
type Bard struct {
	Cookie string
	logger *zerolog.Logger
	answer bardAnswer
}

// New creates a new Bard AI instance. Cookie is the __Secure-1PSID cookie from Google
func New(cookie string, l *zerolog.Logger) *Bard {
	b := &Bard{
		Cookie: cookie,
		logger: l,
	}
	b.answer = bardAnswer{}
	return b
}

// Clear clears the bard answer IDs
func (b *Bard) Clear() {
	b.answer.ChoiceID = ""
	b.answer.ConversationID = ""
	b.answer.ResponseID = ""
}

// Ask generates a Bard AI response and returns it to the user
func (b *Bard) Ask(prompt string) (string, error) {
	// Create a Resty Client
	client := resty.New()

	client.SetLogger(Log{b.logger})
	client.SetDebug(true)

	client.SetHeaders(headers)
	client.SetCookie(&http.Cookie{
		Name:  "__Secure-1PSID",
		Value: b.Cookie,
	})

	// get snim0e value from bard
	client.SetBaseURL("https://bard.google.com")
	client.SetTimeout(5 * time.Second)

	resp, err := client.R().Get("/")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("status code is not 200: %d", resp.StatusCode())
	}

	// req paramters for the actual request
	reqParams := map[string]string{
		"bl":     "boq_assistant-bard-web-server_20230510.09_p1",
		"_reqid": fmt.Sprintf("%d", rand.IntnRange(100000, 999999)),
		"rt":     "c",
	}

	// in response text, the value shows. in python:
	r := regexp.MustCompile(`SNlM0e\":\"(.*?)\"`)

	tmpValues := r.FindStringSubmatch(resp.String())
	if len(tmpValues) < 2 {
		return "", fmt.Errorf("failed to find snim0e value. possibly misconfigured cookies?")
	}
	snim0e := r.FindStringSubmatch(resp.String())[1]

	req := fmt.Sprintf(`[null, "[[\"%s\"], null, [\"%s\", \"%s\", \"%s\"]]"]`,
		//prompt, b.answer.ConversationID, b.answer.ResponseID, b.answer.ChoiceID)
		prompt, b.answer.ConversationID, b.answer.ResponseID, b.answer.ChoiceID)

	reqData := map[string]string{
		"f.req": string(req),
		"at":    snim0e,
	}

	client.SetBaseURL(bardURL)
	client.SetTimeout(60 * time.Second)
	client.SetFormData(reqData)
	client.SetQueryParams(reqParams)
	client.SetDoNotParseResponse(true)
	resp, err = client.R().Post("")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		// curl, _ := http2curl.GetCurlCommand(resp.Request.EnableTrace().RawRequest)
		// fmt.Println(curl)
		return "", fmt.Errorf("status code is not 200: %d", resp.StatusCode())
	}

	// this is the Go version
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.RawResponse.Body)

	respLines := strings.Split(buf.String(), "\n")
	respJSON := respLines[3]

	var fullRes [][]interface{}
	err = json.Unmarshal([]byte(respJSON), &fullRes)
	if err != nil {
		return "", err
	}

	// get the main answer
	err = json.Unmarshal([]byte(fullRes[0][2].(string)), &fullRes)
	if err != nil {
		return "", err
	}

	b.answer.Content = fullRes[0][0].(string)
	b.answer.ConversationID = fullRes[1][0].(string)
	b.answer.ResponseID = fullRes[1][1].(string)

	for _, v := range fullRes[4] {
		choices := v.([]interface{})
		b.answer.ChoiceID = choices[0].(string)
		break
	}

	return b.answer.Content, nil
}
