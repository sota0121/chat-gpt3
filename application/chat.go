package application

type ChatService interface {
}

func NewChatService() ChatService {
	return &chatService{}
}

type chatService struct {
}
