package service

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/anatoliin/url-shortener/internal/model"
)

type URLService struct {
	repo model.URLRepository
}

func NewURLService(repo model.URLRepository) *URLService {
	return &URLService{repo}
}

func (s *URLService) GetShorURL(originalURL string) (string, error) {
	for range 10 {
		short, err := GenerateString(8)
		if err != nil {
			return "", err
		}
		err = s.repo.Save(short, originalURL)
		if err == nil {
			return short, nil
		}
	}
	return "", errors.New("internal error")
}

func (s *URLService) GetOriginalURL(shortURL string) (string, error) {
	return s.repo.FindByID(shortURL)
}

// алгоритм генерации от ИИ
const codeAlphabet = "23456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

func GenerateString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("n must be > 0")
	}

	alphabetLen := big.NewInt(int64(len(codeAlphabet)))
	limit := new(big.Int).Mul(alphabetLen, big.NewInt(256/int64(len(codeAlphabet))))

	buf := make([]byte, length)
	out := make([]byte, 0, length)

	for len(out) < length {
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
		for _, b := range buf {
			if len(out) >= length {
				break
			}
			bi := big.NewInt(int64(b))
			if bi.Cmp(limit) < 0 {
				idx := new(big.Int).Mod(bi, alphabetLen).Int64()
				out = append(out, codeAlphabet[idx])
			}
		}
	}
	return string(out), nil
}
