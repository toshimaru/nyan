package styles

import (
	"testing"

	"github.com/alecthomas/chroma/v2"
	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	names := Names()

	assert.Equal(t, []string{"abap", "dracula", "emacs", "monokai", "monokailight", "pygments", "solarized-dark", "solarized-light", "swapoff", "vim"}, names)
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
