package main

import (
	"appconfig"
	"fmt"
	"requestserver"
)

func main() {
	config := appconfig.GetConfig()

	fmt.Printf("%+v\n", config)

	fmt.Println("HTTP interface for api...")
	requestserver.RunServer()
}
