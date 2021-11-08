package gen

import "embed"

var (
	//go:embed **/*.swagger.json
	SwaggerFS embed.FS
)
