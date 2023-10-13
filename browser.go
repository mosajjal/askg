package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func getCookiesFromBrowser() {
	cookieJar := make(map[string]string)

	opts := append(chromedp.DefaultExecAllocatorOptions[3:],
		// if user-data-dir is set, chrome won't load the default profile,
		// even if it's set to the directory where the default profile is stored.
		// set it to empty to prevent chromedp from setting it to a temp directory.
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("profile-directory", "Default"),
		// BUG: this needs to be cross-platform
		chromedp.UserDataDir(`/home/ali/.config/google-chrome`),
		//chromedp.UserDataDir(os.Getenv("CHROME_USER_DIR")),
		// in headless mode, chrome won't load the default profile.
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-extensions", false),
	)
	ctx_, cancel_ := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel_()

	ctx, cancel := chromedp.NewContext(
		ctx_,
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	err := chromedp.Run(ctx,
		// BUG: if there's a window open, this fails. It should be able to use the open window.
		chromedp.Navigate("https://bard.google.com"),
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
		}))

	if err != nil {
		logger.Fatal().Msgf("Failed to retrieve cookies: %v", err)
	}

	fmt.Printf("Cookie jar: %v\n", cookieJar)

}
