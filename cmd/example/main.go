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

	envs, err := client.GetAll()
	if err != nil {
		log.Fatalf("failed to get env: %v", err)
	}

	fmt.Println(envs)

	env, err := client.GetEnvironmentByName("development")
	if err != nil {
		log.Fatalf("failed to get env: %v", err)
	}

	fmt.Println(env)
}
