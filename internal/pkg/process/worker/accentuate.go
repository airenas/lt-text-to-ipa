package worker

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
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
func NewAccentuator(urlStr string) (process.Processor, error) {
	res := &accentuator{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *accentuator) Process(data *process.Data) error {
	inData := mapAccentInput(data)
	if len(inData) > 0 {

		var output []accentOutputElement
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

type accentOutputElement struct {
	Accent []accent `json:"accent"`
	Word   string   `json:"word"`
	Error  string   `json:"error"`
}

type accent struct {
	MF       string                  `json:"mf"`
	Mi       string                  `json:"mi"`
	MiVdu    string                  `json:"mi_vdu"`
	Mih      string                  `json:"mih"`
	Error    string                  `json:"error"`
	Variants []process.AccentVariant `json:"variants"`
}

func mapAccentInput(data *process.Data) []string {
	res := []string{}
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.IsWord() && w.UserTranscription == "" {
			res = append(res, w.Tagged.Word)
		}
	}
	return res
}

func mapAccentOutput(data *process.Data, out []accentOutputElement) error {
	i := 0
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.IsWord() && w.UserTranscription == "" {
			if len(out) <= i {
				return errors.New("Wrong accent result")
			}
			err := setAccent(w, out[i])
			if err != nil {
				return err
			}
			i++
		}
	}
	return nil
}

func setAccent(w *process.ProcessedWord, out accentOutputElement) error {
	if out.Error != "" {
		if len(w.Tagged.Word) >= 50 {
			goapp.Log.Error(out.Error)
			return utils.NewErrWordTooLong(w.Tagged.Word)
		}
		return errors.Errorf("Accent error for '%s'('%s'): %s", w.Tagged.Word, out.Word, out.Error)
	}
	if w.Tagged.Word != out.Word {
		return errors.Errorf("Words do not match '%s' vs '%s'", w.Tagged.Word, out.Word)
	}
	w.AccentVariant = findBestAccentVariant(out.Accent, w.Tagged.Mi, w.Tagged.Lemma)
	return nil
}

func findBestAccentVariant(acc []accent, mi string, lema string) *process.AccentVariant {
	find := func(fa func(a *accent) bool, fv func(v *process.AccentVariant) bool) *process.AccentVariant {
		for _, a := range acc {
			if fa(&a) {
				for _, v := range a.Variants {
					if fv(&v) {
						return &v
					}
				}
			}
		}
		return nil
	}
	fIsAccent := func(v *process.AccentVariant) bool { return v.Accent > 0 }

	if res := find(func(a *accent) bool { return a.Error == "" && a.MiVdu == mi && a.MF == lema }, fIsAccent); res != nil {
		return res
	}

	if res := find(func(a *accent) bool { return a.Error == "" && a.MiVdu == mi }, fIsAccent); res != nil {
		return res
	}
	// no mi filter
	if res := find(func(a *accent) bool { return a.Error == "" }, fIsAccent); res != nil {
		return res
	}
	//no filter
	return find(func(a *accent) bool { return true }, func(v *process.AccentVariant) bool { return true })
}