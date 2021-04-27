package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/petergtz/pegomock"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewClitics(t *testing.T) {
	initTestJSON(t)
	pr, err := NewClitics("http://server")
	assert.Nil(t, err)
	assert.NotNil(t, pr)
}

func TestNewClitics_Fails(t *testing.T) {
	initTestJSON(t)
	pr, err := NewClitics("")
	assert.NotNil(t, err)
	assert.Nil(t, pr)
}

func TestInvokeClitics(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewClitics("http://server")
	assert.NotNil(t, pr)
	pr.(*cliticDetector).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &extapi.AccentVariant{Accent: 103}})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]cliticsOutput) = []cliticsOutput{{ID: 0,
				Type: "CLITIC"}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.NotNil(t, d.Words[0].Clitic)
}

func TestInvokeClitics_Fail(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewClitics("http://server")
	assert.NotNil(t, pr)
	pr.(*cliticDetector).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &extapi.AccentVariant{Accent: 103}})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).ThenReturn(
		errors.New("olia err"))
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestMapCliticsInput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "olia", Type: process.Word,
		Mi: "mi", Lemma: "lemma"},
		AccentVariant: &extapi.AccentVariant{Accent: 103, Syll: "o-lia"}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &extapi.AccentVariant{Accent: 103}})
	inp, err := mapCliticsInput(d)
	assert.Nil(t, err)
	assert.Equal(t, "olia", inp[0].String)
	assert.Equal(t, "mi", inp[0].Mi)
	assert.Equal(t, "lemma", inp[0].Lemma)
	assert.Equal(t, 0, inp[0].ID)
}

func TestMapCliticsOutput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("olia")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})

	output := []cliticsOutput{{ID: 0, Type: "CLITIC", Accent: 103},
		{ID: 1, Type: "PHRASE", AccentType: "NONE"}}

	err := mapCliticsOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, "CLITIC", d.Words[0].Clitic.Type)
	assert.Equal(t, "NONE", d.Words[1].Clitic.AccentedType)
}

func TestToType(t *testing.T) {
	assert.Equal(t, "OTHER", toType(process.OtherWord))
	assert.Equal(t, "SPACE", toType(process.Space))
	assert.Equal(t, "OTHER", toType(process.SentenceEnd))
	assert.Equal(t, "WORD", toType(process.Word))
	assert.Equal(t, "OTHER", toType(process.Separator))
}
