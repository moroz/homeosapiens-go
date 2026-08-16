package handlers

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"maragu.dev/gomponents"
)

func wrapRender(template gomponents.Node, output http.ResponseWriter) error {
	var result bytes.Buffer

	if err := template.Render(&result); err != nil {
		slog.Error("Template render failed", "err", err)
		return echo.NewHTTPError(500, "Internal Server Error")
	}

	output.Header().Set("Content-Type", "text/html; charset=utf-8")
	output.WriteHeader(200)
	_, err := io.Copy(output, &result)
	return err
}
