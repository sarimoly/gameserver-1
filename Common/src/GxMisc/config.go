package GxMisc

import (
	"fmt"
	. "github.com/bitly/go-simplejson"
	"os"
)

var Config *Json

func LoadConfig(filename string) {
	buf := make([]byte, 1024)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()
	f.Read(buf)
	//
	Config, _ = NewJson(buf)
}
