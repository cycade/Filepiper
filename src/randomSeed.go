package main

import (
    "fmt"
    "time"
    "math/rand"
)
// 随机生成四位提取码
const seed = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
func ecodeGenerator() string {
	randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	extractCode := ""
	for i := 0; i < 4; i++ {
		extractCode += string(seed[randomNumber.Intn(61)])
	}
	return extractCode
}

func main() {
	fmt.Println(ecodeGenerator())
}