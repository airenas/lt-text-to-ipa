package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCleaner(t *testing.T) {
	pr := NewCleaner()
	assert.NotNil(t, pr)
}

func TestInvokeCleaner(t *testing.T) {
	pr := NewCleaner()
	d := newTestData()
	d.OriginalText = "olia\r    kélias"
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, "olia\n    kelias", d.Text)
}
