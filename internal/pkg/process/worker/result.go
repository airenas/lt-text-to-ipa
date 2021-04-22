package worker

import (
	"fmt"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/pkg/errors"
)

type ipaTypeEnum int

const (
	None ipaTypeEnum = iota + 1
	WordOne
	WordMultiple
	Sep
)

var ipaTypeStringEnum = map[ipaTypeEnum]string{None: "NONE", WordOne: "ONE", WordMultiple: "MULTIPLE", Sep: "SEP"}

func ipaToString(t ipaTypeEnum) string {
	return ipaTypeStringEnum[t]
}

type resultMaker struct {
}

//NewResultMaker makes preocessor for result construction
func NewResultMaker() process.Processor {
	res := &resultMaker{}
	return res
}

func (p *resultMaker) Process(data *process.Data) error {
	var err error
	data.Result, err = mapResult(data)
	if err != nil {
		return errors.Wrap(err, "can't prepare results")
	}
	return nil
}

func mapResult(data *process.Data) ([]*api.ResultWord, error) {
	res := make([]*api.ResultWord, 0)
	for _, w := range data.Words {
		tgw := w.Tagged
		if w.Tagged.Type == process.Word {
			res = append(res, &api.ResultWord{Type: "WORD", String: w.Tagged.String, IPA: w.IPA,
				IPAType: getIPAWordType(w)})
		} else if w.Tagged.Type == process.OtherWord {
			res = append(res, &api.ResultWord{Type: "WORD", String: w.Tagged.String, IPAType: ipaToString(None)})
		} else if w.Tagged.Type == process.SentenceEnd {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "//",
				IPAType: ipaToString(Sep)})
		} else if w.Tagged.Type == process.Separator && tgw.String == "," {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "/",
				IPAType: ipaToString(Sep)})
		} else if w.Tagged.Type == process.Separator && tgw.String == "\n" {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "\n",
				IPAType: ipaToString(Sep)})
		} else {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: w.Tagged.String,
				IPA: fmt.Sprintf("%*s", len(tgw.String), " "), IPAType: ipaToString(None)})
		}
	}
	return res, nil
}

func getIPAWordType(w *process.ProcessedWord) string {
	if w.AccentCount > 1 || w.TranscriptionCount > 1 {
		return ipaToString(WordMultiple)
	}
	return ipaToString(WordOne)
}
