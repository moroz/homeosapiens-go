package main

import "github.com/google/uuid"

type ThumbnailProps struct {
	PP     string `json:"pp"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Host   string `json:"host"`
	Locale string `json:"locale"` // "en" | "pl"
}

type VideoItem struct {
	Locale         string
	ID             uuid.UUID
	Title          string
	Host           string
	ProfilePicture string
	RecordedOn     string
}
