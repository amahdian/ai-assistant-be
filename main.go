package main

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/amahdian/ai-assistant-be/global/env"
	"github.com/amahdian/ai-assistant-be/server"
)

func init() {
	// ensure server is always working in UTC
	time.Local = time.UTC
}

//	@title			My App
//	@version		2.0
//	@description	Swagger documentation for the My App's RESTful API.

//	@query.collection.format	multi

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
// @x-extension-openapi		{"example": "value on a json format"}
func main() {
	envs, err := env.Load("")
	if err != nil {
		log.Fatalf("failed to load env variables: %v", err)
	}

	s, err := server.NewServer(envs)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
