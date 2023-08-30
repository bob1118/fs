package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

// string compare case-insensitivity
func IsEqual(s string, d string) (b bool) {
	return strings.EqualFold(s, d)
}

// MakeA1Hash function.
func MakeA1Hash(in string) (s string) {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// UUIDFormat function.
//
// cfeda278-7476-441a-8e2b-c4fe7ec8eb28;Outbound Call;8002
//
// cfeda278-7476-441a-8e2b-c4fe7ec8eb28
func UUIDFormat(in string) string {
	var out string
	var tmp []string
	if len(in) > 36 {
		tmp = strings.Split(in, ";")
		out = tmp[0]
	} else {
		out = in
	}
	return out
}

// UUIDFormatEx function.
//
// ARRAY::16c63f9f-f5d7-4b81-8fe6-cbb5cd8f182a;Outbound Call;8003|:a3d41ee6-ee09-47d2-aa9c-2998df858948;Outbound Call;8002
//
// return slice {"16c63f9f-f5d7-4b81-8fe6-cbb5cd8f182a","a3d41ee6-ee09-47d2-aa9c-2998df858948"}
func UUIDFormatEx(in string) []string {
	var tmp, out []string
	if len(in) > 36 {
		ss := strings.TrimPrefix(in, "ARRAY::")
		tmp = strings.Split(ss, "|:")
		for index, value := range tmp {
			out[index] = UUIDFormat(value)
		}
	} else {
		out[0] = in
	}
	return out
}
