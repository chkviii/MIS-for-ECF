package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"mypage-backend/internal/config"
)

func Html_Path(raw string) string {
	cfg := config.Load()
	return path.Join(cfg.Html_Path, raw)
}

// General function to parse JSON request body into a struct
func ParseJSONBody(r *http.Request, v interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("content-type must be application/json")
	}

	// // 限制请求体大小（可选）
	// r.Body = http.MaxBytesReader(nil, r.Body, 1048576) // 1MB

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(v)
}

func AtoUi(s string) (i uint) {

	t, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}

	i = uint(t) // possible overflow for 32-bit systems, but unlikely in practice

	return i
}

func UitoA(i uint) (s string) {
	s = strconv.FormatUint(uint64(i), 10)
	return s
}
