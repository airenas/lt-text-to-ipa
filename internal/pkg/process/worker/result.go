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
	SepClitic
)

var ipaTypeStringEnum = map[ipaTypeEnum]string{None: "NONE", WordOne: "ONE", WordMultiple: "MULTIPLE",
	SepClitic: "SEP_CLITIC"}

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
	l := len(data.Words)

	for i := 0; i < l; i++ {
		w := data.Words[i]
		if w.Clitic != nil && w.Clitic.Type == "PHRASE" {
			rw, ni := collectPhrase(data.Words, i)
			rw.Info = makeInfo(rw.IPA, []string{"morfologinė samplaika"})
			rw.Info.Word = rw.String
			res = append(res, rw)
			i = ni
			continue
		}
		tgw := w.Tagged

		// drop last value if only one word wanted
		if i == (l-1) && len(res) == 1 && tgw.Type == process.SentenceEnd {
			break
		}

		if w.Tagged.Type == process.Word {
			rw := &api.ResultWord{Type: "WORD", String: tgw.String, IPA: w.IPA,
				Trans:   getTrans(data.ReturnSAMPA, w.Transcription),
				IPAType: getIPAWordType(w)}
			if w.Clitic != nil && w.Clitic.Type == "CLITIC" && w.Clitic.AccentedType == "NONE" {
				rw.Info = makeInfo(rw.IPA, w.Mihs)
				rw.Info.Word = tgw.String
				rw.IPAType = "ONE"
			}
			res = append(res, rw)
		} else if w.Tagged.Type == process.OtherWord {
			res = append(res, &api.ResultWord{Type: "WORD", String: tgw.String, IPA: tgw.String, IPAType: ipaToString(None)})
		} else if w.Tagged.Type == process.SentenceEnd {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "//",
				IPAType: ipaToString(None)})
		} else if w.Tagged.Type == process.Space && betweenClitics(data.Words, i) {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "‿",
				IPAType: ipaToString(SepClitic)})
		} else if w.Tagged.Type == process.Separator && tgw.String == "," {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "/",
				IPAType: ipaToString(None)})
		} else if w.Tagged.Type == process.Separator && tgw.String == "\n" {
			res = append(res, &api.ResultWord{Type: "SEPARATOR", String: tgw.String, IPA: "\n",
				IPAType: ipaToString(None)})
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

func betweenClitics(words []*process.ProcessedWord, at int) bool {
	if at == 0 || at == len(words)-1 {
		return false
	}
	return words[at-1].Clitic != nil && words[at+1].Clitic != nil && (words[at-1].Clitic.Pos+1) == words[at+1].Clitic.Pos
}

func collectPhrase(words []*process.ProcessedWord, at int) (*api.ResultWord, int) {
	l := len(words)
	res := &api.ResultWord{Type: "PHRASE", IPAType: "ONE"}
	li := at
	ni := 0
	for i := at; i < l; i++ {
		w := words[i]
		tgw := w.Tagged
		if w.Tagged.Type == process.Word {
			if w.Clitic != nil && w.Clitic.Type == "PHRASE" && w.Clitic.Pos == ni {
				if ni == 0 {
					res.IPA = w.IPA
					res.String = tgw.String
				} else {
					res.IPA = res.IPA + "‿" + w.IPA
					res.String = res.String + " " + tgw.String
				}
				ni++
				li = i
			} else {
				break
			}
		}
	}
	return res, li
}

func makeInfo(ipa string, mis []string) *api.WordInfo {
	res := &api.WordInfo{}
	rmis := make([]api.MIInfo, len(mis))
	for i, mi := range mis {
		rmis[i].MI = mi
	}
	res.Transcriptions = []*api.Transcription{{IPAs: []string{ipa}, Information: rmis}}
	return res
}

func getTrans(use bool, tr string) string {
	if use {
		return tr
	}
	return ""
}
