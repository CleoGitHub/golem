package stringtool

import "github.com/segmentio/go-snakecase"

func SnakeCase(s string) string {
	return snakecase.Snakecase(s)
}
