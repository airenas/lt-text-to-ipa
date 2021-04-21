package worker

import (
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/tts-line/pkg/clean"
)

type cleaner struct {
}

//NewCleaner makes new text clean processor
func NewCleaner() process.Processor {
	res := &cleaner{}
	return res
}

func (p *cleaner) Process(data *process.Data) error {
	data.Text = clean.ChangeSymbols(clean.DropEmojis(data.OriginalText))
	return nil
}
