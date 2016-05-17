package main

import (
    "crypto/md5"
    "fmt"
    "io"
    "os"
)

func main() {
    file, err := os.Open("245.png")
    if err != nil {
        panic(err)
    }
 
    h := md5.New()
    _, err = io.Copy(h, file)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%x\n", h.Sum(nil))
    // output: 43c6359298645ded23f3c2ee44acf564
}