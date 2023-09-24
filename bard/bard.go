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
	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/util/rand"
)

var headers map[string]string = map[string]string{
	"Host":          "bard.google.com",
	"X-Same-Domain": "1",
	"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	"Content-Type":  "application/x-www-form-urlencoded;charset=utf-8",
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
	Cookie1PSID   string
	Cookie1PSIDTS string
	Cookie1PSIDCC string
	logger        *zerolog.Logger
	answer        bardAnswer

	// Timeout in seconds
	TimeoutSnim0e int
	TimeoutQuery  int
}

// New creates a new Bard AI instance. Cookie is the __Secure-1PSID cookie from Google
func New(cookie1psid, cookie1psidts, cookie1psidcc string, l *zerolog.Logger) *Bard {
	b := &Bard{
		Cookie1PSID:   cookie1psid,
		Cookie1PSIDTS: cookie1psidts,
		Cookie1PSIDCC: cookie1psidcc,
		logger:        l,
		TimeoutSnim0e: 5,
		TimeoutQuery:  60,
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
	client.SetCookies([]*http.Cookie{
		{
			Name:  "__Secure-1PSID",
			Value: b.Cookie1PSID,
		}, {
			Name:  "__Secure-1PSIDCC",
			Value: b.Cookie1PSIDCC,
		}, {
			Name:  "__Secure-1PSIDTS",
			Value: b.Cookie1PSIDTS,
		},
	},
	)

	// get snim0e value from bard
	client.SetBaseURL("https://bard.google.com")
	if b.TimeoutSnim0e > 0 {
		client.SetTimeout(time.Duration(b.TimeoutSnim0e) * time.Second)
	}

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

	var sessionStruct = []interface{}{
		[]string{prompt},
		nil,
		[]string{
			b.answer.ConversationID,
			b.answer.ResponseID,
			b.answer.ChoiceID,
		},
	}

	ls_byte, err := json.Marshal(sessionStruct)
	if err != nil {
		return "", err
	}

	var reqStruct = []interface{}{
		nil,
		string(ls_byte),
	}

	req, err := json.Marshal(reqStruct)
	if err != nil {
		return "", err
	}

	reqData := map[string]string{
		"f.req": string(req),
		"at":    snim0e,
	}

	client.SetBaseURL(bardURL)
	if b.TimeoutQuery > 0 {
		client.SetTimeout(time.Duration(b.TimeoutQuery) * time.Second)
	}
	client.SetFormData(reqData)
	client.SetJSONEscapeHTML(false)
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
	res, ok := fullRes[0][2].(string)
	if !ok {
		return "", fmt.Errorf("failed to get answer from bard")
	}

	b.answer.ConversationID = gjson.Get(res, "1.0").String()
	b.answer.ResponseID = gjson.Get(res, "1.1").String()
	choices := gjson.Get(res, "4").Array()
	b.answer.ChoiceID = choices[0].Array()[0].String()
	b.answer.Content = choices[0].Array()[1].Array()[0].String()

	return b.answer.Content, nil
}
