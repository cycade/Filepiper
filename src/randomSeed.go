package main

import (
    "fmt"
    "time"
    "math/rand"
)
// 随机生成四位提取码
func ecodeProducer() string {
	var seed = [62]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	extractCode := ""
	for i := 0; i < 4; i++ {
		extractCode += seed[randomNumber.Intn(61)]
	}
	return extractCode
}

func ecodeAdder (ecode string)

func main() {
	fmt.Println(ecodeProducer())
}