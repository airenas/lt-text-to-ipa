package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/extapi"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"

	"github.com/stretchr/testify/assert"
)

func TestNewResultMaker(t *testing.T) {
	initTestJSON(t)
	pr := NewResultMaker()
	assert.NotNil(t, pr)
}

func TestInvokeResultMaker(t *testing.T) {
	pr := NewResultMaker()
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"),
		AccentVariant: &extapi.AccentVariant{Accent: 103}, IPA: "w o r d"})
	err := pr.Process(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(d.Result)) {
		assert.Equal(t, "w o r d", d.Result[0].IPA)
	}
}

func TestMapPhrase(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word"), IPA: "w o r d"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word2"), IPA: "w o r d 2"})
	d.Words[0].Clitic = &process.Clitic{Type: "PHRASE", Pos: 0}
	d.Words[2].Clitic = &process.Clitic{Type: "PHRASE", Pos: 1}
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res)) {
		assert.Equal(t, "w o r d‿w o r d 2", res[0].IPA)
		assert.Equal(t, "PHRASE", res[0].Type)
		assert.Equal(t, "ONE", res[0].IPAType)
		assert.Equal(t, "morfologinė samplaika", res[0].Info.Transcriptions[0].Information[0].MI)
	}
}

func TestMapClitic(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("word2"), IPA: "w o r d 2"})
	d.Words[0].Clitic = &process.Clitic{Type: "CLITIC", Pos: 0, AccentedType: "NONE"}
	d.Words[0].Mihs = []string{"mi info"}
	d.Words[2].Clitic = &process.Clitic{Type: "CLITIC", Pos: 1}
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 3, len(res)) {
		assert.Equal(t, "n e", res[0].IPA)
		assert.Equal(t, "WORD", res[0].Type)
		assert.Equal(t, "ONE", res[0].IPAType)
		assert.Equal(t, "mi info", res[0].Info.Transcriptions[0].Information[0].MI)
		assert.Equal(t, "‿", res[1].IPA)
		assert.Equal(t, "w o r d 2", res[2].IPA)
	}
}

func TestMapPause(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 3, len(res)) {
		assert.Equal(t, "n e", res[0].IPA)
		assert.Equal(t, " \u2016 ", res[1].IPA)
		assert.Equal(t, "n e", res[2].IPA)
	}
}

func TestMapPauseWithSep(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 4, len(res)) {
		assert.Equal(t, "n e", res[0].IPA)
		assert.Equal(t, " ", res[1].IPA)
		assert.Equal(t, "\u2016 ", res[2].IPA)
		assert.Equal(t, "n e", res[3].IPA)
	}
}

func TestMapPauseWithSep2(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 5, len(res)) {
		assert.Equal(t, "n e", res[0].IPA)
		assert.Equal(t, " ", res[1].IPA)
		assert.Equal(t, "\u2016", res[2].IPA)
		assert.Equal(t, " ", res[3].IPA)
		assert.Equal(t, "n e", res[4].IPA)
	}
}

func TestMapCommaWithSep2(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 5, len(res)) {
		assert.Equal(t, "n e", res[0].IPA)
		assert.Equal(t, " ", res[1].IPA)
		assert.Equal(t, "\u007C", res[2].IPA)
		assert.Equal(t, " ", res[3].IPA)
		assert.Equal(t, "n e", res[4].IPA)
	}
}

func TestMapCommaWithSep3(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSpace(" ")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 5, len(res)) {
		assert.Equal(t, " ", res[0].IPA)
		assert.Equal(t, "\u007C ", res[1].IPA)
		assert.Equal(t, "\u007C ", res[2].IPA)
		assert.Equal(t, "\u2016 ", res[3].IPA)
		assert.Equal(t, "\u2016", res[4].IPA)
	}
}

func TestMapSep(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep("\n")})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res)) {
		assert.Equal(t, "\n", res[0].IPA)
		assert.Equal(t, "NONE", res[0].IPAType)
		assert.Equal(t, "SEPARATOR", res[0].Type)
	}
}

func TestMapSepComma(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTSep(",")})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res)) {
		assert.Equal(t, "|", res[0].IPA)
		assert.Equal(t, "NONE", res[0].IPAType)
		assert.Equal(t, "SEPARATOR", res[0].Type)
	}
}

func TestMapSepSentence(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res)) {
		assert.Equal(t, "\u2016", res[0].IPA)
		assert.Equal(t, "NONE", res[0].IPAType)
		assert.Equal(t, "SEPARATOR", res[0].Type)
	}
}

func TestMapOtherWord(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("8num")})
	d.Words[0].Tagged.Type = process.OtherWord
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res)) {
		assert.Equal(t, "8num", res[0].IPA)
		assert.Equal(t, "NONE", res[0].IPAType)
		assert.Equal(t, "WORD", res[0].Type)
	}
}

func TestMapOneWord_DropSentenceEnd(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res)) {
		assert.Equal(t, "n e", res[0].IPA)
		assert.Equal(t, "ONE", res[0].IPAType)
		assert.Equal(t, "WORD", res[0].Type)
	}
}

func TestMapSeveralWords(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: newTestTWord("ne"), IPA: "n e"})
	d.Words = append(d.Words, &process.ProcessedWord{Tagged: process.TaggedWord{Type: process.SentenceEnd}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(res))
}
