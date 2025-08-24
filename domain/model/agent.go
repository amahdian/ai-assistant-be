package model

type Agent struct {
	Name         string
	Description  string
	SystemPrompt string
}

var AllAgents []*Agent

var DefaultAgent = &Agent{
	Name:         "Default",
	SystemPrompt: "You are a helpful assistant for the AI-Assistant App. You are powered by a sophisticated AI model.",
}
