package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/petergtz/pegomock"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewTranscriber(t *testing.T) {
	initTestJSON(t)
	pr, err := NewTranscriber("http://server")
	assert.Nil(t, err)
	assert.NotNil(t, pr)
}

func TestNewTranscriber_Fails(t *testing.T) {
	initTestJSON(t)
	pr, err := NewTranscriber("")
	assert.NotNil(t, err)
	assert.Nil(t, pr)
}

func TestInvokeTranscriber(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewTranscriber("http://server")
	assert.NotNil(t, pr)
	pr.(*transcriber).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &process.AccentVariant{Accent: 103}})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]transOutput) = []transOutput{{Word: "word",
				Transcription: []trans{{Transcription: "w o r d"}}}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, "w o r d", d.Words[0].Transcription)
}

func TestInvokeTranscriber_FailInput(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewTranscriber("http://server")
	assert.NotNil(t, pr)
	pr.(*transcriber).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: nil})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).ThenReturn(errors.New("haha"))
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestInvokeTranscriber_Fail(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewTranscriber("http://server")
	assert.NotNil(t, pr)
	pr.(*transcriber).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &process.AccentVariant{Accent: 103}})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).ThenReturn(errors.New("haha"))
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestInvokeTranscriber_NoData(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewTranscriber("http://server")
	assert.NotNil(t, pr)
	d := newTestData()
	err := pr.Process(d)
	assert.Nil(t, err)
}

func TestInvokeTranscriber_FailOutput(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewTranscriber("http://server")
	assert.NotNil(t, pr)
	pr.(*transcriber).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &process.AccentVariant{Accent: 103}})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]transOutput) = []transOutput{}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestMapTransInput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("olia"),
		AccentVariant: &process.AccentVariant{Accent: 103, Syll: "o-lia"}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &process.AccentVariant{Accent: 103}})
	inp, err := mapTransInput(d)
	assert.Nil(t, err)
	assert.Equal(t, "olia", inp[0].Word)
	assert.Equal(t, 103, inp[0].Acc)
	assert.Equal(t, "o-lia", inp[0].Syll)
	assert.Equal(t, "word", inp[1].Word)
	assert.Equal(t, 103, inp[1].Acc)
	assert.Equal(t, "", inp[1].Syll)
}

// func TestMapTransInput_RC(t *testing.T) {
// 	d := newTestData()
// 	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v1"),
// 		AccentVariant: &process.AccentVariant{Accent: 103, Syll: "o-lia"}})
// 	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
// 		AccentVariant: &process.AccentVariant{Accent: 103}})
// 	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word1"),
// 		AccentVariant: &process.AccentVariant{Accent: 103}})
// 	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
// 	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word2"),
// 		AccentVariant: &process.AccentVariant{Accent: 103}})
// 	inp, err := mapTransInput(d)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "word", inp[0].Rc)
// 	assert.Equal(t, "word1", inp[1].Rc)
// 	assert.Equal(t, "", inp[2].Rc)
// }

func TestMapTransInput_Space(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("olia"), AccentVariant: &process.AccentVariant{Accent: 103, Syll: "o-lia"}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	inp, err := mapTransInput(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(inp)) {
		assert.Equal(t, "olia", inp[0].Word)
	}
}

func TestMapTransOutput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("olia")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})

	output := []transOutput{{Word: "olia", Transcription: []trans{{Transcription: "trans"}}},
		{Word: "word", Transcription: []trans{{Transcription: "trans1"}}}}

	err := mapTransOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, "trans", d.Words[0].Transcription)
	assert.Equal(t, "trans1", d.Words[1].Transcription)
}

func TestMapTransOutput_Sep(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("olia")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})

	output := []transOutput{{Word: "olia", Transcription: []trans{{Transcription: "trans"}}},
		{Word: "word", Transcription: []trans{{Transcription: "trans1"}}}}

	err := mapTransOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, "trans", d.Words[1].Transcription)
	assert.Equal(t, "trans1", d.Words[3].Transcription)
}

func TestMapTransOutput_DropQMark(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})

	output := []transOutput{{Word: "word", Transcription: []trans{{Transcription: "tran? - s1?"}}}}

	err := mapTransOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, "tran - s1", d.Words[0].Transcription)
}

func TestMapTransOutput_FailLen(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v1")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})

	output := []transOutput{{Word: "olia", Transcription: []trans{{Transcription: "trans"}}}}

	err := mapTransOutput(d, output)
	assert.NotNil(t, err)
}

func TestMapTransOutput_FailWord(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v1")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})

	output := []transOutput{{Word: "olia1", Transcription: []trans{{Transcription: "trans"}}},
		{Word: "word", Transcription: []trans{{Transcription: "trans1"}}}}

	err := mapTransOutput(d, output)
	assert.NotNil(t, err)
}

func TestMapTransOutput_FailError(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v1")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})

	output := []transOutput{{Word: "olia", Transcription: []trans{{Transcription: "trans"}}},
		{Word: "word", Error: "err"}}

	err := mapTransOutput(d, output)
	assert.NotNil(t, err)
}

func newTestTWord(w string) process.TaggedWord {
	return process.TaggedWord{String: w, Type: process.Word}
}

func newTestTSep(s string) process.TaggedWord {
	return process.TaggedWord{String: s, Type: process.Separator}
}

func newTestTSpace(s string) process.TaggedWord {
	return process.TaggedWord{String: s, Type: process.Space}
}
