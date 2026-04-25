package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
)

type ThumbnailProps struct {
	PP     string `json:"pp"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Host   string `json:"host"`
	Locale string `json:"locale"` // "en" | "pl"
}

const script = `
(async () => {
	const mod = await import("http://localhost:5174/src/main.ts");
	return mod.renderThumbnail;
})()
`

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: new(false),
	})
	if err != nil {
		log.Fatal(err)
	}

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  640,
			Height: 9.0 / 16 * 640,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = page.Goto("http://localhost:5174", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		log.Fatal(err)
	}

	page.On("console", func(msg playwright.ConsoleMessage) {
		fmt.Println("BROWSER:", msg.Text())
	})

	handle, err := page.EvaluateHandle(script)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Dispose()

	props, err := json.Marshal(&ThumbnailProps{
		PP:     "https://d3n1g0yg3ja4p3.cloudfront.net/019beef9-ad4c-736f-9bb0-965b59ca21ae.png",
		Title:  "What prevents me from moving on?",
		Date:   "2026-02-08",
		Host:   "Dr Asher Shaikh",
		Locale: "en",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = handle.Evaluate(`(fn, props) => fn(props)`, string(props))
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

	if err := pw.Stop(); err != nil {
		log.Fatal(err)
	}
}
