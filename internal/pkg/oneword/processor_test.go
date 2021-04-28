package oneword

import (
	"errors"
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/test/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	processorMock *procMock
	worker        *MainWorker
)

func initTest(t *testing.T) {
	mocks.AttachMockToTest(t)
	processorMock = &procMock{f: func(d *Data) error { return nil }}
	worker = &MainWorker{}
	worker.Add(processorMock)
}

func TestWork(t *testing.T) {
	initTest(t)
	processorMock.f = func(d *Data) error {
		assert.Equal(t, "olia", d.Word)
		d.Result = &api.WordInfo{}
		return nil
	}
	res, err := worker.Process("olia", false)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestWork_Fails(t *testing.T) {
	initTest(t)
	processorMock.f = func(d *Data) error {
		return errors.New("olia")
	}
	res, err := worker.Process("olia", false)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestWork_Several(t *testing.T) {
	initTest(t)
	processorMock.f = func(d *Data) error {
		return nil
	}
	processorMock1 := &procMock{f: func(d *Data) error {
		d.Result = &api.WordInfo{Word: "olia"}
		return nil
	}}
	worker.Add(processorMock1)
}

func TestWork_StopProcess(t *testing.T) {
	initTest(t)
	processorMock.f = func(d *Data) error {
		return errors.New("olia")
	}
	processorMock1 := &procMock{f: func(d *Data) error {
		assert.Fail(t, "Unexpected call")
		return nil
	}}
	worker.Add(processorMock1)
	res, err := worker.Process("olia", false)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

type procMock struct {
	f func(res *Data) error
}

func (pr *procMock) Process(d *Data) error {
	return pr.f(d)
}
