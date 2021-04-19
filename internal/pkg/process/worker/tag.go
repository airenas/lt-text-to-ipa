package worker

import (
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/utils"
	"github.com/pkg/errors"
)

//HTTPInvoker makes http call
type HTTPInvoker interface {
	InvokeText(string, interface{}) error
}

type tagger struct {
	httpWrap HTTPInvoker
}

//NewTagger creates new processor
func NewTagger(urlStr string) (process.Processor, error) {
	res := &tagger{}
	var err error
	res.httpWrap, err = utils.NewHTTWrap(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't init http client")
	}
	return res, nil
}

func (p *tagger) Process(data *process.Data) error {
	var output []*TaggedWord
	err := p.httpWrap.InvokeText(data.OriginalText, &output)
	if err != nil {
		return err
	}
	data.Words = mapTagResult(output)
	return nil
}

//TaggedWord - tagger's result
type TaggedWord struct {
	Type   string
	String string
	Mi     string
	Lemma  string
}

func mapTagResult(tags []*TaggedWord) []*process.ProcessedWord {
	res := make([]*process.ProcessedWord, 0)
	for _, t := range tags {
		pw := process.ProcessedWord{Tagged: mapTag(t)}
		res = append(res, &pw)
	}
	return res
}

func mapTag(tag *TaggedWord) process.TaggedWord {
	res := process.TaggedWord{}
	if tag.Type == "SEPARATOR" {
		res.Separator = tag.String
	} else if tag.Type == "SENTENCE_END" {
		res.SentenceEnd = true
	} else if tag.Type == "WORD" || tag.Type == "NUMBER" {
		res.Word = tag.String
		res.Lemma = tag.Lemma
		res.Mi = tag.Mi
	} else if tag.Type == "SPACE" {
		res.Space = true
	}
	return res
}
