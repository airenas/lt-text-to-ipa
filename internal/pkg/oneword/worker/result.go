package worker

import (
	"sort"
	"strings"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/pkg/errors"
)

type resultMaker struct {
}

//NewResultMaker makes preocessor for result construction
func NewResultMaker() oneword.Processor {
	res := &resultMaker{}
	return res
}

func (p *resultMaker) Process(data *oneword.Data) error {
	var err error
	data.Result, err = mapResult(data)
	if err != nil {
		return errors.Wrap(err, "can't prepare results")
	}
	return nil
}

func mapResult(data *oneword.Data) (*api.WordInfo, error) {
	res := &api.WordInfo{}
	kRes := make(map[string]*api.Transcription)
	res.Word = data.Word
	meanings := collectMeanings(data.Words)
	for _, m := range meanings {
		for _, w := range data.Words {
			if w.Meaning == m {
				k := key(w)
				rt, ok := kRes[k]
				if !ok {
					rt = &api.Transcription{}
					res.Transcriptions = append(res.Transcriptions, rt)
					rt.Information = []api.MIInfo{{MI: w.MI, Meaning: w.Meaning}}
					kRes[k] = rt
				}
				rt.IPAs = appendStrings(rt.IPAs, w.IPAs)
				if (data.ReturnSAMPA) {
					rt.Trans = appendStrings(rt.Trans, w.Transcriptions)
				}
			}
		}
	}
	res.Transcriptions = joinIPAs(res.Transcriptions)
	return res, nil
}

func joinIPAs(data []*api.Transcription) []*api.Transcription {
	kRes := make(map[string]*api.Transcription)
	res := make([]*api.Transcription, 0)
	for _, t := range data {
		k := keyIPA(t.IPAs)
		rt, ok := kRes[k]
		if !ok {
			res = append(res, t)
			kRes[k] = t
		} else {
			rt.Information = append(rt.Information, t.Information...)
		}
	}
	return res
}

func collectMeanings(words []*oneword.WorkingWord) []string {
	res := make([]string, 0)
	add := make(map[string]bool)
	for _, w := range words {
		if _, ok := add[w.Meaning]; !ok {
			res = append(res, w.Meaning)
			add[w.Meaning] = true
		}
	}
	return res
}

func appendStrings(s1, s2 []string) []string {
	add := make(map[string]bool)
	for _, w := range s1 {
		add[w] = true
	}
	res := s1
	for _, w := range s2 {
		if _, ok := add[w]; !ok {
			res = append(res, w)
			add[w] = true
		}
	}
	return res
}

func key(w *oneword.WorkingWord) string {
	return w.Meaning + ":" + w.MI
}

func keyIPA(s []string) string {
	cs := make([]string, len(s))
	copy(cs, s)
	sort.Strings(cs)
	res := strings.Builder{}
	for _, st := range cs {
		res.WriteString(st)
	}
	return res.String()
}
