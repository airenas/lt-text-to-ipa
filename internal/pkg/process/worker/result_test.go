package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"

	"github.com/stretchr/testify/assert"
)

func TestNewResultMaker(t *testing.T) {
	initTestJSON(t)
	pr := NewResultMaker()
	assert.NotNil(t, pr)
}

func TestInvokeResultMaker(t *testing.T) {
	pr := NewResultMaker()
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &process.AccentVariant{Accent: 103}, Transcription: "w o r d"})
	err := pr.Process(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(d.Result)) {
		assert.Equal(t, "w o r d", d.Result[0].IPA)
	}
}
