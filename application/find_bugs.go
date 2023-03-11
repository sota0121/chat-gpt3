package application

type FindBugService interface {
}

func NewFindBugService() FindBugService {
	return &findBugService{}
}

type findBugService struct {
}
