package model

type URL struct {
	ID          string
	OriginalURL string
}

type URLRepository interface {
	FindByID(id string) (string, error)
	Save(id, url string) error
}
