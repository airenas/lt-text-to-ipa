package worker

import (
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/pkg/errors"
)

type resultMaker struct {
}

//NewResultMaker makes preocessor for result construction
func NewResultMaker() process.Processor {
	res := &resultMaker{}
	return res
}

func (p *resultMaker) Process(data *process.Data) error {
	var err error
	data.Result, err = mapResult(data)
	if err != nil {
		return errors.Wrap(err, "can't prepare results")
	}
	return nil
}

func mapResult(data *process.Data) ([]*api.ResultWord, error) {
	res := make([]*api.ResultWord, 0)
	for _, w := range data.Words {
		if w.Tagged.IsWord() {
			res = append(res, &api.ResultWord{Type: "WORD", String: w.Tagged.Word, IPA: w.Transcription, IPAType: ""})
		} else {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: w.Tagged.Separator, IPA: w.Tagged.Separator, IPAType: ""})
		}
	}
	return res, nil
}
