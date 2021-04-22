package worker

import (
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

type transcriber struct {
	httpWrap HTTPInvokerJSON
}

//NewTranscriber creates new processor
func NewTranscriber(urlStr string) (process.Processor, error) {
	res := &transcriber{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *transcriber) Process(data *process.Data) error {
	inData, err := mapTransInput(data)
	if err != nil {
		return err
	}
	if len(inData) > 0 {
		var output []transOutput
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

type transInput struct {
	Word string `json:"word"`
	Syll string `json:"syll"`
	User string `json:"user"`
	Ml   string `json:"ml"`
	Rc   string `json:"rc"`
	Acc  int    `json:"acc"`
}

type transOutput struct {
	Transcription []trans `json:"transcription"`
	Word          string  `json:"word"`
	Error         string  `json:"error"`
}

type trans struct {
	Transcription string `json:"transcription"`
}

func mapTransInput(data *process.Data) ([]*transInput, error) {
	res := []*transInput{}
	var pr *transInput
	_ = pr
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type != process.Word {
			pr = nil
		} else {
			ti := &transInput{}
			tword := transWord(w)
			ti.Word = tword
			if w.AccentVariant == nil {
				return nil, errors.New("No accent variant for " + tword)
			}
			ti.Acc = w.AccentVariant.Accent
			ti.Syll = w.AccentVariant.Syll
			ti.Ml = w.AccentVariant.Ml
			// if pr != nil {
			// 	pr.Rc = tword
			// }
			res = append(res, ti)
			pr = ti
		}
	}
	return res, nil
}

func transWord(w *process.ProcessedWord) string {
	return w.Tagged.String
}

func mapTransOutput(data *process.Data, out []transOutput) error {
	i := 0
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type == process.Word {
			if len(out) <= i {
				return errors.New("Wrong transcribe result")
			}
			err := setTrans(w, out[i])
			if err != nil {
				return err
			}
			i++
		}
	}
	return nil
}

func setTrans(w *process.ProcessedWord, out transOutput) error {
	if out.Error != "" {
		return errors.Errorf("Transcriber error for '%s'('%s'): %s", transWord(w), out.Word, out.Error)
	}
	if transWord(w) != out.Word {
		return errors.Errorf("Words do not match (transcriber) '%s' vs '%s'", transWord(w), out.Word)
	}
	w.TranscriptionCount = len(out.Transcription)
	for _, t := range out.Transcription {
		if t.Transcription != "" {
			w.Transcription = dropQMarks(t.Transcription)
			return nil
		}
	}
	return nil
}

func dropQMarks(s string) string {
	return strings.ReplaceAll(s, "?", "")
}
