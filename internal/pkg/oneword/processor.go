package oneword

import (
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
)

// Processor does one specific step
type Processor interface {
	Process(*Data) error
}

// MainWorker does transcription work
type MainWorker struct {
	processors []Processor
}

//Work is main method
func (mw *MainWorker) Process(input string) (*api.WordInfo, error) {
	data := &Data{}
	data.Word = input
	err := mw.processAll(data)
	if err != nil {
		return nil, err
	}
	return data.Result, nil
}

//Add adds a processor to the end
func (mw *MainWorker) Add(pr Processor) {
	mw.processors = append(mw.processors, pr)
}

func (mw *MainWorker) processAll(data *Data) error {
	for _, pr := range mw.processors {
		err := pr.Process(data)
		if err != nil {
			return err
		}
	}
	return nil
}
