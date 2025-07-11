package util

import (
	"path"

	"mypage-backend/internal/config"
)

func Html_Path(raw string) string {
	cfg := config.Load()
	return path.Join(cfg.Html_Path, raw)
}
