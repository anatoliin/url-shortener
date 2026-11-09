package model

type URL struct {
	ID          string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLRepository interface {
	FindByID(id string) (string, error)
	Save(id, url string) error
}
