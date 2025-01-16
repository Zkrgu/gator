package main

import (
	"fmt"
	"github.com/zkrgu/gator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err == nil {
		fmt.Printf("%v", conf)
	} else {
		fmt.Printf("%v", err)
	}
}
