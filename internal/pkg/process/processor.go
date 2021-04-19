package process

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
func (mw *MainWorker) Process(input string) ([]*api.ResultWord, error) {
	data := &Data{}
	data.OriginalText = input
	err := mw.processAll(data)
	if err != nil {
		return nil, err
	}
	return mapResult(data)
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

func mapResult(data *Data) ([]*api.ResultWord, error) {
	res := make([]*api.ResultWord, 0)
	return res, nil
}
