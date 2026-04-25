package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
)

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: new(false)})
	if err != nil {
		log.Fatal(err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	_, err = page.Goto("https://moroz.dev")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if err := os.MkdirAll("./screenshots", 0o755); err != nil {
		log.Fatal(err)
	}

	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: new(fmt.Sprintf("screenshots/%s.png", uuid.Must(uuid.NewV7()))),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := browser.Close(); err != nil {
		log.Fatal(err)
	}
}
