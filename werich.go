package werich

import (
	"github.com/russross/blackfriday"
)

// Run convert md to html
func Run(input []byte) []byte {
	return blackfriday.Run(input)
}
