package main

import (
	"appconfig"
	"fmt"
	"requestserver"
)

func main() {
	config := appconfig.GetConfig()

	fmt.Printf("%+v", config)

	fmt.Println("HTTP interface for api...")
	go requestserver.RunServer()
}
