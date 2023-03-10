package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

//string compare case-insensitivity
func IsEqual(s string, d string) (b bool) {
	return strings.EqualFold(s, d)
}

//MakeA1Hash function.
func MakeA1Hash(in string) (s string) {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//UUIDFormat function.
//
//cfeda278-7476-441a-8e2b-c4fe7ec8eb28;Outbound Call;8002
//
//cfeda278-7476-441a-8e2b-c4fe7ec8eb28
func UUIDFormat(in string) (s string) {
	var out string
	var tmp []string
	if len(in) > 36 {
		tmp = strings.Split(in, ";")
		out = tmp[0]
	} else {
		out = in
	}
	return strings.TrimSpace(out)
}
