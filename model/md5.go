package model
import(
    "io"

    "crypto/md5"
	//"crypto/sha512"
	"encoding/hex"
)

func Md5(code string) string{
    hasher := md5.New()
	_, _ = io.WriteString(hasher, code)
	return hex.EncodeToString(hasher.Sum(nil))
}