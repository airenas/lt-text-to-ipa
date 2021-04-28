package worker

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

type toIPA struct {
	httpWrap HTTPInvokerJSON
}

//NewAccentuator creates new processor
func NewToIPA(urlStr string) (oneword.Processor, error) {
	res := &toIPA{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *toIPA) Process(data *oneword.Data) error {
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

func mapToIPAInput(data *oneword.Data) []*extapi.IPAInput {
	rs := make(map[string]bool)
	for _, w := range data.Words {
		for _, t := range w.Transcriptions {
			rs[t] = true
		}
	}
	res := make([]*extapi.IPAInput, 0)
	for k := range rs {
		res = append(res, &extapi.IPAInput{Transcription: k})
	}
	return res
}

func mapIPAOutput(data *oneword.Data, out []extapi.IPAOutput) error {
	rs := make(map[string]string)
	for _, ipa := range out {
		rs[ipa.Transcription] = ipa.IPA

	}
	for _, w := range data.Words {
		ipas := make([]string, 0)
		for _, t := range w.Transcriptions {
			ipas = append(ipas, rs[t])
		}
		w.IPAs = ipas
	}
	return nil
}
