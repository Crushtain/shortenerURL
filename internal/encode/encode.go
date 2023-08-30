package encode

import "github.com/jxskiss/base62"

func Encode(body []byte) []byte {
	short := base62.Encode(body)
	return short
}
