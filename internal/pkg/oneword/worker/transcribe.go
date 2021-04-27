package worker

import (
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

type transcriber struct {
	httpWrap HTTPInvokerJSON
}

//NewTranscriber creates new processor
func NewTranscriber(urlStr string) (oneword.Processor, error) {
	res := &transcriber{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *transcriber) Process(data *oneword.Data) error {
	inData, err := mapTransInput(data)
	if err != nil {
		return err
	}
	if len(inData) > 0 {
		var output []extapi.TransOutput
		err := p.httpWrap.InvokeJSON(inData, &output)
		if err != nil {
			return err
		}
		err = mapTransOutput(data, output)
		if err != nil {
			return err
		}
	} else {
		goapp.Log.Debug("Skip transcriber - no data in")
	}
	return nil
}

func mapTransInput(data *oneword.Data) ([]*extapi.TransInput, error) {
	res := []*extapi.TransInput{}
	for _, w := range data.Words {
		ti := &extapi.TransInput{}
		ti.Word = data.Word
		ti.Acc = w.Accent
		ti.Syll = w.Syll
		ti.Ml = w.Lemma
		res = append(res, ti)
	}
	return res, nil
}

func mapTransOutput(data *oneword.Data, out []extapi.TransOutput) error {
	i := 0
	for _, w := range data.Words {
		if len(out) <= i {
			return errors.New("wrong transcribe result")
		}
		if out[0].Error != "" {
			return errors.Errorf("transcriber error for '%s'('%s'): %s", data.Word, out[0].Word, out[0].Error)
		}
		w.Transcriptions = make([]string, 0)
		for _, t := range out[i].Transcription {
			if t.Transcription != "" {
				w.Transcriptions = append(w.Transcriptions, dropQMarks(t.Transcription))
			}
		}
		i++
	}
	return nil
}

func dropQMarks(s string) string {
	return strings.ReplaceAll(s, "?", "")
}
