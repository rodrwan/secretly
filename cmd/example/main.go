package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rodrwan/secretly/pkg/secretly"
)

func main() {
	client := secretly.New(
		secretly.WithBaseURL("http://localhost:8080/api/v1"),
	)

	env, err := client.GetEnv()
	if err != nil {
		log.Fatalf("failed to get env: %v", err)
	}

	fmt.Println(env)

	client.LoadToEnvironment()

	fmt.Println(os.Getenv("BEBESITA"))
}
