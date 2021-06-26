package assets

import "embed"

var (
	//go:embed project/*.* project project/internal/model/*.* project/internal/dto/*.*
	ProjectFS embed.FS
)
