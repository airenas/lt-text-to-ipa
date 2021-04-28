package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
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
	d.Words = append(d.Words, &oneword.WorkingWord{Accent: 102, Syll: "o-lia"})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]extapi.TransOutput) = []extapi.TransOutput{{Word: "word",
				Transcription: []extapi.Trans{{Transcription: "w o r d"}}}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, []string{"w o r d"}, d.Words[0].Transcriptions)
}

func TestInvokeTranscriber_Fail(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewTranscriber("http://server")
	assert.NotNil(t, pr)
	pr.(*transcriber).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Accent: 102, Syll: "o-lia"})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).ThenReturn(errors.New("haha"))
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestMapTransInput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Accent: 102, Syll: "o-lia"})
	d.Words = append(d.Words, &oneword.WorkingWord{Accent: 103, Syll: "o-lia2"})
	inp, err := mapTransInput(d)
	assert.Nil(t, err)
	assert.Equal(t, "olia", inp[0].Word)
	assert.Equal(t, 102, inp[0].Acc)
	assert.Equal(t, "o-lia", inp[0].Syll)
	assert.Equal(t, "olia", inp[1].Word)
	assert.Equal(t, 103, inp[1].Acc)
	assert.Equal(t, "o-lia2", inp[1].Syll)
}

func TestMapTransOutput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Accent: 102, Syll: "o-lia"})
	d.Words = append(d.Words, &oneword.WorkingWord{Accent: 103, Syll: "o-lia2"})

	output := []extapi.TransOutput{{Word: "olia", Transcription: []extapi.Trans{{Transcription: "tr1"}, {Transcription: "tr2"}}},
		{Word: "olia", Transcription: []extapi.Trans{{Transcription: "trans1 - ?a"}}}}

	err := mapTransOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, []string{"tr1", "tr2"}, d.Words[0].Transcriptions)
	assert.Equal(t, []string{"trans1 - a"}, d.Words[1].Transcriptions)
}
