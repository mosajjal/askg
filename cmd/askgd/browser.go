package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func getChromeDefaultProfilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	var chromeProfilePath string
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		chromeProfilePath = filepath.Join(localAppData, `Google\Chrome\User Data\`)
	case "darwin":
		chromeProfilePath = filepath.Join(usr.HomeDir, `Library/Application Support/Google/Chrome/`)
	case "linux":
		chromeProfilePath = filepath.Join(usr.HomeDir, `.config/google-chrome/`)
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return chromeProfilePath, nil
}

func getCookiesFromBrowser(profile string) map[string]string {
	if profile == "" {
		profile = "Default"
	}
	cookieJar := make(map[string]string)
	dir, _ := getChromeDefaultProfilePath()

	opts := append(chromedp.DefaultExecAllocatorOptions[3:],
		// if user-data-dir is set, chrome won't load the default profile,
		// even if it's set to the directory where the default profile is stored.
		// set it to empty to prevent chromedp from setting it to a temp directory.
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("profile-directory", profile),
		chromedp.UserDataDir(dir),
		//chromedp.UserDataDir(os.Getenv("CHROME_USER_DIR")),
		// in headless mode, chrome won't load the default profile.
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-extensions", false),
	)
	ctx_, cancel_ := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel_()

	ctx, cancel := chromedp.NewContext(
		ctx_,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://gemini.google.com/app"),
		chromedp.WaitVisible("div[ng-non-bindable]", chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				fmt.Println("get cookies error", err)
				return err
			}
			for _, cookie := range cookies {
				cookieJar[cookie.Name] = cookie.Value
			}
			return nil
		})); err != nil {
		logger.Fatal().Msgf("Failed to retrieve cookies: %v", err)
	}

	return cookieJar
}
