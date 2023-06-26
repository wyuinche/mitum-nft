package digest

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

var JSON = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            false,
	ValidateJsonRawMessage: false,
}.Froze()

func Marshal(v interface{}) ([]byte, error) {
	return JSON.Marshal(v)
}

func Unmarshal(b []byte, v interface{}) error {
	return JSON.Unmarshal(b, v)
}

func parseLimitQuery(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(-1)
	}
	return n
}

func parseStringQuery(s string) string {
	return strings.TrimSpace(s)
}

func stringOffsetQuery(offset string) string {
	return fmt.Sprintf("offset=%s", offset)
}

func parseBoolQuery(s string) bool {
	return s == "1"
}

func stringBoolQuery(key string, v bool) string { // nolint:unparam
	if v {
		return fmt.Sprintf("%s=1", key)
	}

	return ""
}

func addQueryValue(b, s string) string {
	if len(s) < 1 {
		return b
	}

	if !strings.Contains(b, "?") {
		return b + "?" + s
	}

	return b + "&" + s
}
