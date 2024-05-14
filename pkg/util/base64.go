package util

import "encoding/base64"

func ToImageBase64(data []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
}
