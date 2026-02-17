package services

import (
	"fmt"

	"github.com/koki-develop/go-fzf"
)

type selectConfig interface {
	Theme() string
	Layout() string
}

type SelectService struct {
	config selectConfig
}

func NewSelectService(config selectConfig) SelectService {
	return SelectService{config: config}
}

func themeOptions(theme string) []fzf.Option {
	switch theme {
	case "default":
		return nil
	case "bw":
		return []fzf.Option{
			fzf.WithStyles(
				fzf.WithStyleCursor(fzf.Style{Faint: true}),
				fzf.WithStyleCursorLine(fzf.Style{Bold: true}),
				fzf.WithStyleMatches(fzf.Style{Bold: true}),
				fzf.WithStyleSelectedPrefix(fzf.Style{Faint: true}),
			),
		}
	case "minimal":
		return []fzf.Option{
			fzf.WithStyles(
				fzf.WithStyleCursor(fzf.Style{ForegroundColor: "#666666"}),
				fzf.WithStyleMatches(fzf.Style{ForegroundColor: "#888888"}),
				fzf.WithStyleSelectedPrefix(fzf.Style{ForegroundColor: "#666666"}),
			),
		}
	default:
		return nil
	}
}

func layoutToInputPosition(layout string) fzf.InputPosition {
	if layout == "reverse" {
		return fzf.InputPositionTop
	}
	return fzf.InputPositionBottom
}

func (s SelectService) Select(options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	fzfOpts := []fzf.Option{
		fzf.WithInputPosition(layoutToInputPosition(s.config.Layout())),
	}
	fzfOpts = append(fzfOpts, themeOptions(s.config.Theme())...)

	f, err := fzf.New(fzfOpts...)
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
