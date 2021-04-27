package worker

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

type cliticDetector struct {
	httpWrap HTTPInvokerJSON
}

//NewClitics creates new processor
func NewClitics(urlStr string) (process.Processor, error) {
	res := &cliticDetector{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *cliticDetector) Process(data *process.Data) error {
	inData, err := mapCliticsInput(data)
	if err != nil {
		return err
	}
	if len(inData) > 0 {
		var output []cliticsOutput
		err := p.httpWrap.InvokeJSON(inData, &output)
		if err != nil {
			return err
		}
		err = mapCliticsOutput(data, output)
		if err != nil {
			return err
		}
	} else {
		goapp.Log.Debug("Skip transcriber - no data in")
	}
	return nil
}

type cliticsInput struct {
	Type   string `json:"type,omitempty"`
	String string `json:"string,omitempty"`
	Mi     string `json:"mi,omitempty"`
	Lemma  string `json:"lemma,omitempty"`
	ID     int    `json:"id,omitempty"`
}

type cliticsOutput struct {
	ID         int    `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	AccentType string `json:"accentType,omitempty"`
	Accent     int    `json:"accent,omitempty"`
	Pos        int    `json:"pos,omitempty"`
}

func mapCliticsInput(data *process.Data) ([]*cliticsInput, error) {
	res := []*cliticsInput{}
	for i, w := range data.Words {
		tgw := w.Tagged
		ci := &cliticsInput{}
		ci.ID = i
		ci.String = transWord(w)
		ci.Lemma = tgw.Lemma
		ci.Mi = tgw.Mi
		ci.Type = toType(w.Tagged.Type)
		res = append(res, ci)
	}
	return res, nil
}

func toType(t process.StringTypeEnum) string {
	if t == process.Word {
		return "WORD"
	}
	if t == process.Space {
		return "SPACE"
	}
	return "OTHER"
}

func mapCliticsOutput(data *process.Data, out []cliticsOutput) error {
	om := make(map[int]cliticsOutput)
	for _, co := range out {
		om[co.ID] = co
	}

	for i, w := range data.Words {
		fco, ok := om[i]
		if ok {
			w.Clitic = &process.Clitic{Accent: fco.Accent, AccentedType: fco.AccentType,
				Type: fco.Type, Pos: fco.Pos}
		}
	}
	return nil
}
