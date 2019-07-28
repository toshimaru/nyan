package styles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	names := Names()

	assert.Equal(t, []string{"dracula", "monokai", "solarized-dark", "swapoff", "vim"}, names)
}
