package main

import (
	"fmt"

	"github.com/ivandersr/products-api-go/configs"
)

func main() {
	conf := configs.LoadConfig(".")
	payload := make(map[string]interface{})
	fmt.Println(conf.TokenAuth.Encode(payload))
}
