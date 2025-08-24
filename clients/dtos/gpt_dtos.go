package dtos

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}
