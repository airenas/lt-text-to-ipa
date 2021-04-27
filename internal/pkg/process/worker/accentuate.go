package worker

import (
	"fmt"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
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

func mapAccentInput(data *process.Data) []string {
	res := []string{}
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type == process.Word {
			res = append(res, w.Tagged.String)
		}
	}
	return res
}

func mapAccentOutput(data *process.Data, out []extapi.AccentOutputElement) error {
	i := 0
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type == process.Word {
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

func setAccent(w *process.ProcessedWord, out extapi.AccentOutputElement) error {
	if out.Error != "" {
		if len(w.Tagged.String) >= 50 {
			goapp.Log.Error(out.Error)
			return utils.NewErrWordTooLong(w.Tagged.String)
		}
		return errors.Errorf("Accent error for '%s'('%s'): %s", w.Tagged.String, out.Word, out.Error)
	}
	if w.Tagged.String != out.Word {
		return errors.Errorf("Words do not match '%s' vs '%s'", w.Tagged.String, out.Word)
	}
	w.AccentVariant = findBestAccentVariant(out.Accent, w.Tagged.Mi, w.Tagged.Lemma)
	w.Mihs = collectMihs(out.Accent)
	w.AccentCount = countVariants(out.Accent)
	return nil
}

func findBestAccentVariant(acc []extapi.Accent, mi string, lema string) *extapi.AccentVariant {
	find := func(fa func(a *extapi.Accent) bool, fv func(v *extapi.AccentVariant) bool) *extapi.AccentVariant {
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
	fIsAccent := func(v *extapi.AccentVariant) bool { return v.Accent > 0 }

	if res := find(func(a *extapi.Accent) bool { return a.Error == "" && a.MiVdu == mi && a.MF == lema }, fIsAccent); res != nil {
		return res
	}

	if res := find(func(a *extapi.Accent) bool { return a.Error == "" && a.MiVdu == mi }, fIsAccent); res != nil {
		return res
	}
	// no mi filter
	if res := find(func(a *extapi.Accent) bool { return a.Error == "" }, fIsAccent); res != nil {
		return res
	}
	//no filter
	return find(func(a *extapi.Accent) bool { return true }, func(v *extapi.AccentVariant) bool { return true })
}

func countVariants(acc []extapi.Accent) int {
	am := make(map[string]bool)
	for _, a := range acc {
		for _, v := range a.Variants {
			if v.Accent > 0 {
				am[fmt.Sprintf("%d,%s,%s", v.Accent, v.Syll, a.MF)] = true
			}
		}
	}
	return len(am)
}

func collectMihs(acc []extapi.Accent) []string {
	res := make([]string, 0)
	for _, a := range acc {
		if a.Mih != "" {
			res = append(res, a.Mih)
		}
	}
	return res
}
