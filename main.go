package main

import "GoAPI/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)

}
