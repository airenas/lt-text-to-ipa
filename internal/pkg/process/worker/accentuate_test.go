package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
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
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]extapi.AccentOutputElement) = []extapi.AccentOutputElement{{Word: "word",
				Accent: []extapi.Accent{{Mi: "mi", Variants: []extapi.AccentVariant{{Accent: 101}}}}}}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.Nil(t, err)
	assert.Equal(t, 101, d.Words[0].AccentVariant.Accent)
}

func TestInvokeAccentuator_FailOutput(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewAccentuator("http://server")
	assert.NotNil(t, pr)
	pr.(*accentuator).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).Then(
		func(params []pegomock.Param) pegomock.ReturnValues {
			*params[1].(*[]extapi.AccentOutputElement) = []extapi.AccentOutputElement{}
			return []pegomock.ReturnValue{nil}
		})
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestInvokeAccentuator_Fail(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewAccentuator("http://server")
	assert.NotNil(t, pr)
	pr.(*accentuator).httpWrap = httpJSONMock
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word")})
	pegomock.When(httpJSONMock.InvokeJSON(pegomock.AnyInterface(), pegomock.AnyInterface())).ThenReturn(errors.New("haha"))
	err := pr.Process(d)
	assert.NotNil(t, err)
}

func TestInvokeAccentuator_NoData(t *testing.T) {
	initTestJSON(t)
	pr, _ := NewAccentuator("http://server")
	assert.NotNil(t, pr)
	d := newTestData()
	err := pr.Process(d)
	assert.Nil(t, err)
}

func TestMapAccInput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v1")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep("!")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v2")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v3", Type: process.OtherWord}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	inp := mapAccentInput(d)
	assert.Equal(t, []string{"v1", "v2"}, inp)
}

func TestMapAccOutput(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v1", Type: process.OtherWord}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep("!")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("v2")})

	output := []extapi.AccentOutputElement{{Word: "v2",
		Accent: []extapi.Accent{{Mi: "mi", Variants: []extapi.AccentVariant{{Accent: 101,
			Syll: "v-1"}}}}}}

	err := mapAccentOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, 101, d.Words[3].AccentVariant.Accent)
	assert.Equal(t, "v-1", d.Words[3].AccentVariant.Syll)
}

func TestMapAccOutput_FindBest(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v1", Type: process.OtherWord}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep("!")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v2", Mi: "mi2", Type: process.Word}})

	output := []extapi.AccentOutputElement{{Word: "v2",
		Accent: []extapi.Accent{{MiVdu: "mi1", Variants: []extapi.AccentVariant{{Accent: 101,
			Syll: "v-1"}}},
			{MiVdu: "mi2", Variants: []extapi.AccentVariant{{Accent: 102,
				Syll: "v-1"}}},
		}}}

	err := mapAccentOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, 102, d.Words[2].AccentVariant.Accent)
}

func TestMapAccOutput_Error(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v1", Mi: "mi2", Type: process.OtherWord}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep("!")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v2", Mi: "mi2", Type: process.Word}})

	output := []extapi.AccentOutputElement{{Word: "v2",
		Accent: []extapi.Accent{{MiVdu: "mi1", Error: "err", Variants: []extapi.AccentVariant{{Accent: 0,
			Syll: "v-1"}}},
			{MiVdu: "mi2", Variants: []extapi.AccentVariant{{Accent: 102,
				Syll: "v-2"}}},
		}}}

	err := mapAccentOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, "v-2", d.Words[2].AccentVariant.Syll)
}

func TestMapAccOutput_WithAccent(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep("!")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v2", Mi: "mi2", Type: process.Word}})

	output := []extapi.AccentOutputElement{{Word: "v2",
		Accent: []extapi.Accent{
			{MiVdu: "mi1", Error: "err", Variants: []extapi.AccentVariant{{Accent: 0, Syll: "v-1"}}},
			{MiVdu: "mi2", Variants: []extapi.AccentVariant{
				{Accent: 0, Syll: "v-2"},
				{Accent: 103, Syll: "v-3"},
			}},
		}}}

	err := mapAccentOutput(d, output)
	assert.Nil(t, err)
	assert.Equal(t, "v-3", d.Words[1].AccentVariant.Syll)
	assert.Equal(t, 103, d.Words[1].AccentVariant.Accent)
}

func TestMapAccOutput_FailError(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v2", Mi: "mi2", Type: process.Word}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v2", Mi: "mi2", Type: process.Word}})

	output := []extapi.AccentOutputElement{{Word: "v2",
		Accent: []extapi.Accent{
			{MiVdu: "mi1", Variants: []extapi.AccentVariant{{Accent: 0, Syll: "v-1"}}},
		}},
		{Word: "v2", Error: "error olia"}}

	err := mapAccentOutput(d, output)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Accent error for 'v2'('v2'): error olia", err.Error())
	}
}

func TestMapAccOutput_FailErrorTooLong(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{String: "v2", Mi: "mi2", Type: process.Word}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{
		String: "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong", Mi: "mi2", Type: process.Word}})

	output := []extapi.AccentOutputElement{{Word: "v2",
		Accent: []extapi.Accent{
			{MiVdu: "mi1", Variants: []extapi.AccentVariant{{Accent: 0, Syll: "v-1"}}},
		}},
		{Word: "v2", Error: "error olia"}}

	err := mapAccentOutput(d, output)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Wrong accent, too long word: ")
	}
}

func TestFindBest_UseLemma(t *testing.T) {
	acc := []extapi.Accent{{MiVdu: "mi2", MF: "lema1", Variants: []extapi.AccentVariant{{Accent: 101}}},
		{MiVdu: "mi2", MF: "lema", Variants: []extapi.AccentVariant{{Accent: 103}}}}
	res := findBestAccentVariant(acc, "mi2", "lema")

	assert.Equal(t, 103, res.Accent)
}

func newTestData() *process.Data {
	return &process.Data{Text: "olia"}
}
