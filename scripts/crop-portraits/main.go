// crop-portraits downloads host profile pictures from the CDN, center-crops
// them to the portrait aspect ratio used in thumbnails (240:328), uploads
// to S3 under a new key, and updates the object_key in the database.
package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	bucket      = "homeosapiens-staging-assets"
	region      = "ap-northeast-1"
	cdnBase     = "https://d3n1g0yg3ja4p3.cloudfront.net"
	targetW     = 480 // 240 SVG units × 2 for retina
	targetH     = 656 // 328 SVG units × 2 for retina
	jpegQuality = 85
)

const query = `
	SELECT a.id, a.object_key
	FROM assets a
	JOIN hosts h ON h.profile_picture_id = a.id
`

const updateQuery = `UPDATE assets SET object_key = $1 WHERE id = $2`

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type row struct {
		id  uuid.UUID
		key string
	}
	var assets []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.key); err != nil {
			log.Fatal(err)
		}
		assets = append(assets, r)
	}

	if len(assets) == 0 {
		log.Println("no host profile pictures found")
		return
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(awsCfg)

	for _, a := range assets {
		newKey, err := processAndUpload(ctx, s3Client, a.key)
		if err != nil {
			log.Printf("SKIP %s (%s): %v", a.key, a.id, err)
			continue
		}
		if _, err := db.Exec(ctx, updateQuery, newKey, a.id); err != nil {
			log.Printf("DB UPDATE FAILED %s → %s: %v", a.key, newKey, err)
			continue
		}
		log.Printf("OK   %s → %s", a.key, newKey)
	}
}

func processAndUpload(ctx context.Context, client *s3.Client, key string) (string, error) {
	url := fmt.Sprintf("%s/%s", cdnBase, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	src, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	cropped := coverCrop(src, targetW, targetH)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, cropped, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return "", fmt.Errorf("encode: %w", err)
	}

	newKey := fmt.Sprintf("portraits/%s.jpg", uuid.Must(uuid.NewV7()))
	b := bucket
	ct := "image/jpeg"
	cc := "public, max-age=31536000, immutable"

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       &b,
		Key:          &newKey,
		ContentType:  &ct,
		CacheControl: &cc,
		Body:         bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", fmt.Errorf("s3 upload: %w", err)
	}

	return newKey, nil
}

// coverCrop scales src to fill targetW×targetH from the center (object-fit: cover).
func coverCrop(src image.Image, targetW, targetH int) image.Image {
	sb := src.Bounds()
	srcW, srcH := sb.Dx(), sb.Dy()

	scaleX := float64(targetW) / float64(srcW)
	scaleY := float64(targetH) / float64(srcH)
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	scaledW := int(float64(srcW)*scale + 0.5)
	scaledH := int(float64(srcH)*scale + 0.5)

	scaled := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))
	draw.BiLinear.Scale(scaled, scaled.Bounds(), src, sb, draw.Src, nil)

	offsetX := (scaledW - targetW) / 2
	offsetY := (scaledH - targetH) / 2

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.Draw(dst, dst.Bounds(), scaled, image.Pt(offsetX, offsetY), draw.Src)
	return dst
}
