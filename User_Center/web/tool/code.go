package tool

import (
	"math/rand"
	"time"
)

// 生成指定位数的验证码数字篇
func Get_Rand_Code(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano())) //seed已经弃用//高并发seed会有问题
	numbers := "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(result)
}
