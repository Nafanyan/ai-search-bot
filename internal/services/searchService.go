package services

type SearchService interface {
	Search(query string) (string, error)
}

type MockSearchService struct{}

func (m *MockSearchService) Search(query string) (string, error) {
	return "Здесь должен быть ответ от AI модели", nil
}
