package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/test/mocks"
	"github.com/petergtz/pegomock"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

var (
	httpJSONMock *mocks.MockHTTPInvokerJSON
)

func initTestJSON(t *testing.T) {
	mocks.AttachMockToTest(t)
	httpJSONMock = mocks.NewMockHTTPInvokerJSON()
}

func TestNewAccentuator(t *testing.T) {
	initTestJSON(t)
	pr, err := NewAccentuator("http://server")
	assert.Nil(t, err)
	assert.NotNil(t, pr)
}

func TestNewAccentuator_Fails(t *testing.T) {
	initTestJSON(t)
	pr, err := NewAccentuator("")
	assert.NotNil(t, err)
	assert.Nil(t, pr)
}

func TestInvokeAccentuator(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewAccentuator("http://server")
	assert.NotNil(t, pr)
	pr.(*accentuator).httpWrap = httpJSONMock
	d := newTestData()
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]extapi.AccentOutputElement) = []extapi.AccentOutputElement{{Word: "olia",
				Accent: []extapi.Accent{{Mi: "mi", Mih: "mih", Variants: []extapi.AccentVariant{{Accent: 101}}}}}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, "mih", d.Words[0].MI)
	assert.Equal(t, 101, d.Words[0].Accent)
}

func TestInvokeAccentuator_Fail(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewAccentuator("http://server")
	assert.NotNil(t, pr)
	pr.(*accentuator).httpWrap = httpJSONMock
	d := newTestData()
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).ThenReturn(errors.New("haha"))
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestMapAccOutput(t *testing.T) {
	d := newTestData()

	output := []extapi.AccentOutputElement{{Word: "olia",
		Accent: []extapi.Accent{{Mih: "mi", Variants: []extapi.AccentVariant{{Accent: 101,
			Syll: "v-1"}}}}}}

	err := mapAccentOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, 101, d.Words[0].Accent)
	assert.Equal(t, "v-1", d.Words[0].Syll)
}

func newTestData() *oneword.Data {
	return &oneword.Data{Word: "olia"}
}
