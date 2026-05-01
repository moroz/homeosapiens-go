package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/playwright-community/playwright-go"
)

const AssetsHost = "http://localhost:5173"

const script = `
(async () => {
	const mod = await import("` + AssetsHost + `/src/main.ts");
	return { renderThumbnail: mod.renderThumbnail, destroyApp: mod.destroyApp };
})()
`

func getVideos(db *pgxpool.Pool, ctx context.Context) ([]*VideoItem, error) {
	rows, err := db.Query(ctx, listVideosQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*VideoItem
	for rows.Next() {
		var i VideoItem
		err := rows.Scan(
			&i.Locale,
			&i.ID,
			&i.Title,
			&i.Host,
			&i.ProfilePicture,
			&i.RecordedOn,
		)
		if err != nil {
			return result, err
		}
		result = append(result, &i)
	}
	return result, nil
}

func initPlaywright() (playwright.Browser, func(), error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: new("/usr/bin/chromium"),
	})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = browser.Close()
		_ = pw.Stop()
	}

	return browser, cleanup, nil
}

func buildProps(v *VideoItem) *ThumbnailProps {
	return &ThumbnailProps{
		PP:     v.ProfilePicture,
		Title:  v.Title,
		Date:   v.RecordedOn,
		Host:   v.Host,
		Locale: v.Locale,
	}
}

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFrom

	videos, err := getVideos(db, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	browser, cleanup, err := initPlaywright()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  640,
			Height: 9.0 / 16 * 640,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = page.Goto(AssetsHost, playwright.PageGotoOptions{
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

	render, err := handle.GetProperty("renderThumbnail")
	if err != nil {
		log.Fatal(err)
	}

	destroy, err := handle.GetProperty("destroyApp")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("./screenshots", 0o755); err != nil {
		log.Fatal(err)
	}

	for _, video := range videos {
		assetID := uuid.Must(uuid.NewV7())

		props, err := json.Marshal(buildProps(video))
		if err != nil {
			log.Fatal(err)
		}

		_, err = render.Evaluate(`(fn, props) => fn(props)`, string(props))
		if err != nil {
			log.Fatal(err)
		}

		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path: new(fmt.Sprintf("screenshots/%s.png", assetID)),
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = destroy.Evaluate("(fn) => fn()")
		if err != nil {
			log.Fatal(err)
		}

		db.Exec(context.Background(), insertAssetQuery, assetID, fmt.Sprintf("images/%s.png", assetID))
		db.Exec(context.Background(), setThumbnailQuery, video.Locale, assetID, video.ID)
	}
}
