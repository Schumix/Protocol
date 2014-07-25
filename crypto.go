package protocol

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func sha1_gen(data string) string {
	chiperer := sha1.New()
	chiperer.Write([]byte(data))
	bs := chiperer.Sum(nil)
	return hex.EncodeToString(bs)
}

func md5_gen(data string) string {
	chiperer := md5.New()
	chiperer.Write([]byte(data))
	bs := chiperer.Sum(nil)
	return hex.EncodeToString(bs)
}
