package swaggerui

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed ui.html
	ui string
)

const (
	base = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@3.51.0"
)

func Serve(prefix string, swaggerFs embed.FS) echo.MiddlewareFunc {

	tp := template.Must(template.New("index").Parse(ui))
	var index bytes.Buffer
	tp.Execute(&index, map[string]interface{}{
		"base": base,
	})

	files := make([]map[string]string, 0, 32)
	fs.WalkDir(swaggerFs, ".", func(src string, f fs.DirEntry, err error) error {
		if f.IsDir() {
			return nil
		}
		if !strings.HasSuffix(src, ".json") {
			return nil
		}
		files = append(files, map[string]string{
			"name": strings.TrimSuffix(src, ".swagger.json"),
			"url":  path.Join(prefix, src),
		})
		return nil
	})

	h := http.StripPrefix(prefix, http.FileServer(http.FS(swaggerFs)))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestPath := strings.TrimPrefix(c.Path(), prefix)
			switch requestPath {
			case "", "index.html":
				tp.Execute(c.Response().Writer, map[string]interface{}{
					"base":  base,
					"files": files,
				})
				return nil
			default:
				if strings.HasSuffix(requestPath, ".json") {
					h.ServeHTTP(c.Response(), c.Request())
					return nil
				}
			}

			return c.NoContent(http.StatusNotFound)
		}
	}
}
