package application

type TestGenService interface {
}

func NewTestGenService() TestGenService {
	return &testGenService{}
}

type testGenService struct {
}
