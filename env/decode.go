package env

import (
	v "github.com/hydrati/plugin-loader/utils/container"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func decodeBytesGB18030(b []byte) (string, error) {
	if b == nil {
		return "", nil
	}

	decoder := simplifiedchinese.GB18030.NewDecoder()
	s, err := decoder.Bytes(b)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func DecodeBytesGB18030(b []byte) v.Result[string, error] {
	return v.Resuify(decodeBytesGB18030(b))
}
