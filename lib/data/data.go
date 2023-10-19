package data

import (
  "embed"
)

//go:embed commands/*.json
var EmbeddedFiles embed.FS

