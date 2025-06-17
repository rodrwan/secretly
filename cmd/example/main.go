package main

import (
	"fmt"
	"log"

	"github.com/rodrwan/secretly/pkg/secretly"
)

func main() {
	client := secretly.New(
		secretly.WithBaseURL("http://localhost:8080"),
	)

	env, err := client.GetEnv()
	if err != nil {
		log.Fatalf("failed to get env: %v", err)
	}

	fmt.Println(env)
}
