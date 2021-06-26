package v1

import "embed"

var (
	//go:embed *.swagger.json
	SwaggerFS embed.FS
)
