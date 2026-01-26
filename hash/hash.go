package hash

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"

	"github.com/hdget/utils"
	"github.com/matoous/go-nanoid/v2"
	"github.com/speps/go-hashids/v2"
)

func HashToUint32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write(utils.StringToBytes(s))
	return h.Sum32()
}

func HashString(s string, length int) string {
	return HashBytes(utils.StringToBytes(s), length)
}

func HashBytes(data []byte, length int) string {
	hashValue := fmt.Sprintf("%x", sha256.Sum256(data))
	hdData := hashids.NewData()
	hdData.MinLength = length
	h, _ := hashids.NewWithData(hdData)
	value, _ := h.EncodeHex(hashValue)
	return value
}

// GenerateRandString 生成随机字符串
func GenerateRandString(n int) string {
	randStr, _ := gonanoid.Generate(hashids.DefaultAlphabet, n)
	return randStr
}
