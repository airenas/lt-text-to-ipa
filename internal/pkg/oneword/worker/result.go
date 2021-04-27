package worker

import (
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/pkg/errors"
)

type resultMaker struct {
}

//NewResultMaker makes preocessor for result construction
func NewResultMaker() oneword.Processor {
	res := &resultMaker{}
	return res
}

func (p *resultMaker) Process(data *oneword.Data) error {
	var err error
	data.Result, err = mapResult(data)
	if err != nil {
		return errors.Wrap(err, "can't prepare results")
	}
	return nil
}

func mapResult(data *oneword.Data) (*api.WordInfo, error) {
	res := &api.WordInfo{}
	res.Word = data.Word
	return res, nil
}
