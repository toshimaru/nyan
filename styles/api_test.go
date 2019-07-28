package styles

import (
	"testing"

	"github.com/alecthomas/chroma"
	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	names := Names()

	assert.Equal(t, []string{"dracula", "monokai", "solarized-dark", "swapoff", "vim"}, names)
}

func TestGetValidStyle(t *testing.T) {
	style := Get("vim")

	assert.Equal(t, "vim", style.Name)
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
