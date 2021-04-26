package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/test/mocks"
	"github.com/labstack/echo/v4"
	"github.com/petergtz/pegomock"
	"github.com/stretchr/testify/assert"
)

var (
	tData               *Data
	tEcho               *echo.Echo
	tResp               *httptest.ResponseRecorder
	mockTranscriber     *mocks.MockTranscriber
	mockWordTranscriber *mocks.MockWordTranscriber
)

func initTest(t *testing.T) {
	mockTranscriber = mocks.NewMockTranscriber()
	mockWordTranscriber = mocks.NewMockWordTranscriber()
	tData = &Data{Transcriber: mockTranscriber, WordTranscriber: mockWordTranscriber}
	tEcho = initRoutes(tData)
	tResp = httptest.NewRecorder()
}

func TestLive(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodGet, "/live", nil)

	tEcho.ServeHTTP(tResp, req)
	assert.Equal(t, http.StatusOK, tResp.Code)
	assert.Equal(t, `{"service":"OK"}`, tResp.Body.String())
}

func TestNotFound(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodGet, "/any", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)

	tEcho.ServeHTTP(tResp, req)

	assert.Equal(t, http.StatusNotFound, tResp.Code)
}

func TestTransciber(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodPost, "/ipa", strings.NewReader("mama"))
	pegomock.When(mockTranscriber.Process(pegomock.AnyString())).ThenReturn([]*api.ResultWord{{Type: "WORD", String: "mama", IPA: "m a m a"}}, nil)

	tEcho.ServeHTTP(tResp, req)

	assert.Equal(t, http.StatusOK, tResp.Code)
	assert.Equal(t, `[{"type":"WORD","string":"mama","ipa":"m a m a"}]`,
		strings.TrimSpace(tResp.Body.String()))

}

func TestTranscriber_Fails(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest("POST", "/ipa", strings.NewReader("mama"))
	pegomock.When(mockTranscriber.Process(pegomock.AnyString())).ThenReturn(nil, errors.New("olia"))

	tEcho.ServeHTTP(tResp, req)

	assert.Equal(t, http.StatusInternalServerError, tResp.Code)
}

func TestWordTransciber(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodGet, "/ipa/mama", nil)
	pegomock.When(mockWordTranscriber.Process(pegomock.AnyString())).
		ThenReturn(&api.WordInfo{Transcriptions: []api.Transcription{{IPAs: []string{"m am a"}}}}, nil)

	tEcho.ServeHTTP(tResp, req)

	assert.Equal(t, http.StatusOK, tResp.Code)
	assert.Equal(t, `{"word":"","transcription":[{"ipas":["m am a"],"information":null}]}`,
		strings.TrimSpace(tResp.Body.String()))

}

func TestWordTransciber_Fails(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest("POST", "/ipa", strings.NewReader("mama"))
	pegomock.When(mockTranscriber.Process(pegomock.AnyString())).ThenReturn(nil, errors.New("olia"))

	tEcho.ServeHTTP(tResp, req)

	assert.Equal(t, http.StatusInternalServerError, tResp.Code)
}
