package model
import (
	"math/rand"
	"strconv"
	"time"
)
// GenerateRandomNumber 生成一个 10 位的随机数字字符串
func GenerateRandomNumber() string {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	// 生成 10 位随机数字
	var number string
	for i := 0; i < 10; i++ {
		digit := rand.Intn(10) // 生成 0 到 9 之间的随机整数
		number += strconv.Itoa(digit)
	}
	return number
}
