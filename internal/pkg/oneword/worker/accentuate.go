package worker

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

//HTTPInvokerJSON invoker for json input
type HTTPInvokerJSON interface {
	InvokeJSON(interface{}, interface{}) error
}

type accentuator struct {
	httpWrap HTTPInvokerJSON
}

//NewAccentuator creates new processor
func NewAccentuator(urlStr string) (oneword.Processor, error) {
	res := &accentuator{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "can't init http client")
	}
	return res, nil
}

func (p *accentuator) Process(data *oneword.Data) error {
	inData := mapAccentInput(data)
	if len(inData) > 0 {

		var output []extapi.AccentOutputElement
		err := p.httpWrap.InvokeJSON(inData, &output)
		if err != nil {
			return err
		}
		err = mapAccentOutput(data, output)
		if err != nil {
			return err
		}
	} else {
		goapp.Log.Debug("Skip accenter - no data in")
	}
	return nil
}

func mapAccentInput(data *oneword.Data) []string {
	return []string{data.Word}
}

func mapAccentOutput(data *oneword.Data, out []extapi.AccentOutputElement) error {
	if len(out) != 1 {
		return errors.New("wrong accent output")
	}
	if data.Word != out[0].Word {
		return errors.Errorf("words do not match '%s' vs '%s'", data.Word, out[0].Word)
	}
	for _, a := range out[0].Accent {
		for _, v := range a.Variants {
			if a.Error == "" && v.Accent > 0 {
				ww := &oneword.WorkingWord{}
				ww.Accent = v.Accent
				ww.Lemma = a.MF
				ww.Meaning = v.Meaning
				ww.Syll = v.Syll
				ww.MI = a.Mih
				data.Words = append(data.Words, ww)
			}
		}
	}
	return nil
}
