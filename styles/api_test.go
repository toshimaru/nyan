package styles

import (
	"testing"

	"github.com/alecthomas/chroma"
	"github.com/stretchr/testify/assert"
)

// func TestNames(t *testing.T) {
// 	names := Names()

// 	assert.Equal(t, []string{"dracula", "swapoff", "solarized-dark", "swapoff", "swapoff"}, names)
// }

func TestGetValidStyle(t *testing.T) {
	style := Get("swapoff")

	assert.Equal(t, "swapoff", style.Name)
}

func TestGetInvalidStyle(t *testing.T) {
	style := Get("invalid-style")

	assert.Equal(t, "swapoff", style.Name)
}

func TestRegister(t *testing.T) {
	nesStyle := chroma.MustNewStyle("newstyle", chroma.StyleEntries{})
	style := Register(nesStyle)

	assert.Equal(t, "newstyle", style.Name)
}
