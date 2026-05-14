package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/playwright-community/playwright-go"
)

const AssetsHost = "http://localhost:5173"
const Bucket = "homeosapiens-staging-assets"

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
		ExecutablePath: new("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
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

func callMagick(args ...string) error {
	cmd := exec.Command("convert", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func processImage(assetID uuid.UUID) error {
	outDir := path.Join("screenshots", assetID.String())
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	srcImg := fmt.Sprintf("screenshots/%s.png", assetID)

	variants := [][]string{
		{srcImg, "-quality", "80", "-strip", path.Join(outDir, "2x.webp")},
		{srcImg, "-quality", "80", "-strip", path.Join(outDir, "2x.png")},
		{srcImg, "-quality", "80", "-resize", "50%", "-strip", path.Join(outDir, "1x.webp")},
		{srcImg, "-quality", "80", "-resize", "50%", "-strip", path.Join(outDir, "1x.png")},
	}

	for _, variant := range variants {
		if err := callMagick(variant...); err != nil {
			return err
		}
	}

	return nil
}

func uploadImages(client *s3.Client, assetID uuid.UUID) error {
	outDir := path.Join("screenshots", assetID.String())
	files, err := os.ReadDir(outDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		contentType := "image/png"
		if strings.HasSuffix(file.Name(), ".webp") {
			contentType = "image/webp"
		}

		reader, err := os.Open(path.Join(outDir, file.Name()))
		if err != nil {
			return err
		}

		params := &s3.PutObjectInput{
			Bucket:       new(Bucket),
			Key:          new(fmt.Sprintf("images/%s/%s", assetID, file.Name())),
			CacheControl: new("public, max-age=31536000, immutable"),
			ContentType:  new(contentType),
			Body:         reader,
		}

		if _, err := client.PutObject(context.Background(), params); err != nil {
			return err
		}

		reader.Close()
	}

	return nil
}

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(awsCfg)
	_ = s3Client

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

	assetIDs := make([]uuid.UUID, len(videos))

	for i, video := range videos {
		assetID := uuid.Must(uuid.NewV7())
		assetIDs[i] = assetID

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

		if err := processImage(assetID); err != nil {
			log.Fatal(err)
		}

		if _, err := db.Exec(context.Background(), insertAssetQuery, assetID); err != nil {
			log.Fatal(err)
		}

		if _, err := db.Exec(context.Background(), setThumbnailQuery, video.Locale, assetID, video.ID); err != nil {
			log.Fatal(err)
		}
	}

	for _, assetID := range assetIDs {
		if err := uploadImages(s3Client, assetID); err != nil {
			log.Fatal(err)
		}
	}
}
