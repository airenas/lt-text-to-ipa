package worker

import (
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
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

func mapTransInput(data *process.Data) ([]*extapi.TransInput, error) {
	res := []*extapi.TransInput{}
	var pr *extapi.TransInput
	_ = pr
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type == process.Word {
			ti := &extapi.TransInput{}
			tword := transWord(w)
			ti.Word = tword
			if w.AccentVariant == nil {
				return nil, errors.New("No accent variant for " + tword)
			}
			ti.Acc = w.AccentVariant.Accent
			ti.Syll = w.AccentVariant.Syll
			ti.Ml = w.AccentVariant.Ml
			if w.Clitic != nil {
				if w.Clitic.AccentedType == "NONE" {
					ti.Acc = 0
				} else if w.Clitic.AccentedType == "STATIC" {
					ti.Acc = w.Clitic.Accent
				}
				if pr != nil && w.Clitic.Pos > 0 {
					pr.Rc = tword
				}
				pr = ti
			} else {
				pr = nil
			}
			res = append(res, ti)
		}
	}
	return res, nil
}

func transWord(w *process.ProcessedWord) string {
	return strings.ToLower(w.Tagged.String)
}

func mapTransOutput(data *process.Data, out []extapi.TransOutput) error {
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

func setTrans(w *process.ProcessedWord, out extapi.TransOutput) error {
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
