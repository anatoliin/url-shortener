package model

type Url struct {
	ID          string
	OriginalUrl string
}

type UrlRepository interface {
	FindById(id string) (string, error)
	Save(id, url string) error
}
