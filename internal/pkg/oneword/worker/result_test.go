package worker

import (
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
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
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "aaa", MI: "mih", IPAs: []string{"ipa"}})
	err := pr.Process(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(d.Result.Transcriptions)) {
		assert.Equal(t, []string{"ipa"}, d.Result.Transcriptions[0].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih", Meaning: "aaa"}}, d.Result.Transcriptions[0].Information)
	}
}

func TestMap_Meanings(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "aaa", MI: "mih", IPAs: []string{"ipa"}})
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih", IPAs: []string{"ipa2"}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(res.Transcriptions)) {
		assert.Equal(t, []string{"ipa"}, res.Transcriptions[0].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih", Meaning: "aaa"}}, res.Transcriptions[0].Information)
		assert.Equal(t, []string{"ipa2"}, res.Transcriptions[1].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih", Meaning: ""}}, res.Transcriptions[1].Information)
	}
}

func TestMap_GroupIPAs(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih", IPAs: []string{"ipa", "ipa2"}})
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih2", IPAs: []string{"ipa2", "ipa"}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res.Transcriptions)) {
		assert.Equal(t, []string{"ipa", "ipa2"}, res.Transcriptions[0].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih", Meaning: ""}, {MI: "mih2", Meaning: ""}}, res.Transcriptions[0].Information)
	}
}

func TestMap_GroupIPAsSeveral(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih", IPAs: []string{"ipa", "ipa2"}})
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih2", IPAs: []string{"ipa2", "ipa"}})
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih3", IPAs: []string{"ipa2"}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(res.Transcriptions)) {
		assert.Equal(t, []string{"ipa", "ipa2"}, res.Transcriptions[0].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih", Meaning: ""}, {MI: "mih2", Meaning: ""}}, res.Transcriptions[0].Information)
		assert.Equal(t, []string{"ipa2"}, res.Transcriptions[1].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih3", Meaning: ""}}, res.Transcriptions[1].Information)
	}
}

func TestMap_JoinIPAs(t *testing.T) {
	d := newTestData()
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih", IPAs: []string{"ipa", "ipa2"}})
	d.Words = append(d.Words, &oneword.WorkingWord{Meaning: "", MI: "mih", IPAs: []string{"ipa2", "ipa3"}})
	res, err := mapResult(d)
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(res.Transcriptions)) {
		assert.Equal(t, []string{"ipa", "ipa2", "ipa3"}, res.Transcriptions[0].IPAs)
		assert.Equal(t, []api.MIInfo{{MI: "mih", Meaning: ""}}, res.Transcriptions[0].Information)
	}
}
