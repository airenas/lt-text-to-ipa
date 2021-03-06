package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/petergtz/pegomock"

	"github.com/stretchr/testify/assert"
)

func TestNewToIPA(t *testing.T) {
	initTestJSON(t)
	pr, err := NewToIPA("http://server")
	assert.Nil(t, err)
	assert.NotNil(t, pr)
}

func TestNewToIPA_Fails(t *testing.T) {
	initTestJSON(t)
	pr, err := NewToIPA("")
	assert.NotNil(t, err)
	assert.Nil(t, pr)
}

func TestInvokeToIPA(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewToIPA("http://server")
	assert.NotNil(t, pr)
	pr.(*toIPA).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Transcription: "w o r d", Tagged: newTestTWord("word")})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]extapi.IPAOutput) = []extapi.IPAOutput{{Transcription: "w o r d", IPA: "w ooo rr d"}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, "w ooo rr d", d.Words[0].IPA)
}

func TestMapIPAOutput_Fail(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Transcription: "w o r d", Tagged: newTestTWord("word")})
	err := mapIPAOutput(d, []extapi.IPAOutput{{Transcription: "w o r d 1", IPA: "w ooo rr d"}})
	assert.NotNil(t, err)
	err = mapIPAOutput(d, []extapi.IPAOutput{})
	assert.NotNil(t, err)
}
