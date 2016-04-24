package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	Age        int
	Cats       []string
	Pi         float64
	Perfection []int
	DOB        time.Time // requires `import time`
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("./tomlData.toml", &conf)
	if err != nil {
		fmt.Println(err)
	}
	// if _, err := toml.Decode(tomlData.toml, &conf); err != nil {
	// 	// handle error
	// 	fmt.Println(err)
	// }
	fmt.Println(conf)
}
