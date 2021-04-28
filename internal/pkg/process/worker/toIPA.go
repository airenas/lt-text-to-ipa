package worker

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

type toIPA struct {
	httpWrap HTTPInvokerJSON
}

//NewAccentuator creates new processor
func NewToIPA(urlStr string) (process.Processor, error) {
	res := &toIPA{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *toIPA) Process(data *process.Data) error {
	inData := mapToIPAInput(data)
	if len(inData) > 0 {

		var output []extapi.IPAOutput
		err := p.httpWrap.InvokeJSON(inData, &output)
		if err != nil {
			return err
		}
		err = mapIPAOutput(data, output)
		if err != nil {
			return err
		}
	} else {
		goapp.Log.Debug("Skip toIPA - no data in")
	}
	return nil
}

func mapToIPAInput(data *process.Data) []*extapi.IPAInput {
	res := make([]*extapi.IPAInput, 0)
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type == process.Word {
			res = append(res, &extapi.IPAInput{Transcription: w.Transcription})
		}
	}
	return res
}

func mapIPAOutput(data *process.Data, out []extapi.IPAOutput) error {
	i := 0
	for _, w := range data.Words {
		tgw := w.Tagged
		if tgw.Type == process.Word {
			if len(out) <= i {
				return errors.New("Wrong IPA result")
			}
			if w.Transcription != out[i].Transcription {
				return errors.Errorf("Transcriptions do not match '%s' vs '%s'", w.Transcription, out[i].Transcription)
			}
			w.IPA = out[i].IPA
			i++
		}
	}
	return nil
}
