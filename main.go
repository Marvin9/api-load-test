package main

import (
	"fmt"

	"github.com/Marvin9/api-load-test/pkg"
)

func main() {
	fmt.Print("\033[H\033[2J")
	session := &pkg.Session{
		TargetEndpoint: "http://127.0.0.1:8000",
		Rate:           100,
		Until:          10,
	}
	metadata := session.GenerateMetadata()
	session.LoadTest(metadata)
	session.Success()
	return
}
