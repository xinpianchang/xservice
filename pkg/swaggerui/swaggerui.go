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
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/v2/pkg/log"
)

var (
	//go:embed ui.html
	ui string
)

const (
	base = "https://unpkg.com/swagger-ui-dist@4.8.1"
)

// Serve create swaggerFS middleware for serve swagger UI
func Serve(prefix string, efs ...embed.FS) echo.MiddlewareFunc {

	tp := template.Must(template.New("index").Parse(ui))
	var index bytes.Buffer
	_ = tp.Execute(&index, map[string]interface{}{
		"base": base,
	})

	files := make([]map[string]string, 0, 128)
	fsm := make(map[string]int, 128)
	hs := make(map[int]http.Handler, len(efs))

	for i, f := range efs {
		_ = fs.WalkDir(f, ".", func(src string, f fs.DirEntry, err error) error {
			if f.IsDir() {
				return nil
			}
			if !strings.HasSuffix(src, ".json") {
				return nil
			}

			name := strings.TrimSuffix(src, ".swagger.json")
			url := path.Join(prefix, src)

			if _, ok := fsm[url]; ok {
				log.Warn("file already exists, will be ignored", zap.String("file", src))
				return nil
			}

			fsm[url] = i

			files = append(files, map[string]string{
				"name": name,
				"url":  url,
			})

			return nil
		})

		hs[i] = http.StripPrefix(prefix, http.FileServer(http.FS(f)))
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			switch strings.TrimPrefix(path, prefix) {
			case "", "index.html":
				_ = tp.Execute(c.Response().Writer, map[string]interface{}{
					"base":  base,
					"files": files,
				})
				return nil
			default:
				if strings.HasSuffix(path, ".json") {
					if i, ok := fsm[path]; ok {
						hs[i].ServeHTTP(c.Response(), c.Request())
						return nil
					}
				}
			}

			return c.NoContent(http.StatusNotFound)
		}
	}
}
