package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
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
	d.Words = append(d.Words, &oneword.WorkingWord{Transcriptions: []string{"t1", "t2"}})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]extapi.IPAOutput) = []extapi.IPAOutput{{Transcription: "t1", IPA: "w ooo rr d"},
				{Transcription: "t2", IPA: "t ooo rr d"}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, []string{"w ooo rr d", "t ooo rr d"}, d.Words[0].IPAs)
}
