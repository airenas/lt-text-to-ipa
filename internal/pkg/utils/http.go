package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

//HTTPWrap for http call
type HTTPWrap struct {
	HTTPClient *http.Client
	URL        string
	Timeout    time.Duration
}

//NewHTTWrap creates new wrapper
func NewHTTWrap(urlStr string) (*HTTPWrap, error) {
	res := &HTTPWrap{}
	var err error
	res.URL, err = checkURL(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't parse url '%s'", urlStr)
	}
	res.Timeout = time.Second * 30
	res.HTTPClient = &http.Client{Transport: newTransport()}
	return res, nil
}

func newTransport() http.RoundTripper {
	res := http.DefaultTransport.(*http.Transport).Clone()
	res.MaxIdleConns = 100
	res.MaxConnsPerHost = 100
	res.MaxIdleConnsPerHost = 50
	return res
}

//InvokeText makes http call with text
func (hw *HTTPWrap) InvokeText(dataIn string, dataOut interface{}) error {
	req, err := http.NewRequest(http.MethodPost, hw.URL, strings.NewReader(dataIn))
	if err != nil {
		return errors.Wrapf(err, "Can't prepare request to '%s'", hw.URL)
	}
	if hw.Timeout > 0 {
		ctx, cancelF := context.WithTimeout(context.Background(), hw.Timeout)
		defer cancelF()
		req = req.WithContext(ctx)
	}
	LogData("Input: ", dataIn)
	req.Header.Set("Content-Type", "text/plain")
	return hw.invoke(req, dataOut)
}

//InvokeJSON makes http call with json
func (hw *HTTPWrap) InvokeJSON(dataIn interface{}, dataOut interface{}) error {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(dataIn)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, hw.URL, b)
	if err != nil {
		return errors.Wrapf(err, "Can't prepare request to '%s'", hw.URL)
	}
	if hw.Timeout > 0 {
		ctx, cancelF := context.WithTimeout(context.Background(), hw.Timeout)
		defer cancelF()
		req = req.WithContext(ctx)
	}
	LogData("Input: ", b.String())
	req.Header.Set("Content-Type", "application/json")
	return hw.invoke(req, dataOut)
}

func (hw *HTTPWrap) invoke(req *http.Request, dataOut interface{}) error {
	LogData("Call : ", hw.URL)
	resp, err := hw.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "can't call '%s'", hw.URL)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 10000))
		_ = resp.Body.Close()
	}()
	if err := goapp.ValidateHTTPResp(resp, 100); err != nil {
		return err
	}
	br, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "can't read body")
	}
	LogData("Output: ", string(br))
	err = json.Unmarshal(br, dataOut)
	if err != nil {
		return errors.Wrap(err, "can't decode response")
	}
	return nil
}

//MaxLogDataSize indicates how many bytes of data to log
var MaxLogDataSize = 100

//LogData logs data to debug
func LogData(st string, data string) {
	goapp.Log.Debugf("%s %s", st, trimString(data, MaxLogDataSize))
}

func trimString(data string, size int) string {
	if len(data) > size {
		return data[0:size] + "..."
	}
	return data
}
