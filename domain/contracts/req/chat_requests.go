package req

type SendMessage struct {
	Message string `json:"message" binding:"required"`
}
