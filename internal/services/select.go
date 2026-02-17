package services

import (
	"fmt"

	"github.com/koki-develop/go-fzf"
)

type SelectService struct{}

func NewSelectService() SelectService {
	return SelectService{}
}

func (s SelectService) Select(options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	f, err := fzf.New()
	if err != nil {
		return "", fmt.Errorf("creating finder: %v", err)
	}

	idxs, err := f.Find(options, func(i int) string { return options[i] })
	if err != nil {
		return "", fmt.Errorf("no selection made")
	}

	if len(idxs) == 0 {
		return "", fmt.Errorf("no selection made")
	}

	return options[idxs[0]], nil
}
