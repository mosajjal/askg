package gemini

import (
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

func (g *Gemini) RotateCookies() {
	client := resty.New().R()
	client.SetHeaders(
		map[string]string{
			"authority":       "accounts.google.com",
			"accept":          "*/*",
			"accept-language": "en-US,en;q=0.9,fa;q=0.8",
			"content-type":    "application/json",
			"origin":          "https://accounts.google.com",
			"referer":         "https://accounts.google.com/RotateCookiesPage",
			"sec-gpc":         "1",
			"user-agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		},
	)
	client.SetDebug(true)
	client.SetContentLength(true)
	client.SetDoNotParseResponse(true)

	for k, v := range g.Cookies {
		client.SetCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	// client.SetBaseURL("https://accounts.google.com")
	client.SetBody(`["658", "-3338205198646068790"]`)
	resp, err := client.Post("https://accounts.google.com/RotateCookies")
	if err != nil {
		g.logger.Fatal().Msgf("failed to rotate cookies: %s", err)
	}
	if resp.StatusCode() != 200 {
		g.logger.Fatal().Msgf("failed to rotate cookies, response code %d. Might be a good idea to run the browser command to set the cookies up before running the daemon", resp.StatusCode())
	}

	// set the response cookies back to the g object
	for _, c := range resp.Cookies() {
		g.Cookies[c.Name] = c.Value
	}

}

// commit the cookies to the config file
func (g *Gemini) CommitCookies(cfgFilePath string) {
	// try to open the cfg file
	cfg, err := os.Create(cfgFilePath)
	if err != nil {
		g.logger.Error().Msgf("failed to open config file: %s", err)
	}

	// write the cookies to the file
	y, _ := yaml.Marshal(g.Cookies)
	_, err = cfg.Write(y)
}
