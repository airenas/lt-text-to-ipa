package worker

import (
	"strings"

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
	err := p.httpWrap.InvokeText(data.Text, &output)
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
	res.String = tag.String
	if tag.Type == "SEPARATOR" {
		res.Type = process.Separator
	} else if tag.Type == "SENTENCE_END" {
		res.Type = process.SentenceEnd
	} else if tag.Type == "WORD" {
		res.Lemma = tag.Lemma
		res.Mi = tag.Mi
		res.Type = detectType(res.Mi)
	} else if tag.Type == "NUMBER" {
		res.Type = process.OtherWord
	} else if tag.Type == "SPACE" {
		res.Type = process.Space
	}
	return res
}

func detectType(mi string) process.StringTypeEnum {
	if mi == "" || strings.HasPrefix(mi, "X") || strings.HasPrefix(mi, "Y") {
		return process.OtherWord
	}
	return process.Word
}
