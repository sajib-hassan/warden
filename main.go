package main

import (
	"github.com/joho/godotenv"

	"github.com/sajib-hassan/warden/cmd"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
