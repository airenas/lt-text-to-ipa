package mocks

import (
	"testing"

	"github.com/petergtz/pegomock"
)

//go:generate pegomock generate --package=mocks --output=transcriber.go -m github.com/airenas/lt-text-to-ipa/internal/pkg/service Transcriber

//go:generate pegomock generate --package=mocks --output=wordTranscriber.go -m github.com/airenas/lt-text-to-ipa/internal/pkg/service WordTranscriber

//go:generate pegomock generate --package=mocks --output=httpInvoker.go -m github.com/airenas/lt-text-to-ipa/internal/pkg/process/worker HTTPInvoker

////go:generate pegomock generate --package=mocks --output=httpInvokerJSOM.go -m github.com/airenas/lt-text-to-ipa/internal/pkg/process/worker HTTPInvokerJSON

//AttachMockToTest register pegomock verification to be passed to testing engine
func AttachMockToTest(t *testing.T) {
	pegomock.RegisterMockFailHandler(handleByTest(t))
}

func handleByTest(t *testing.T) pegomock.FailHandler {
	return func(message string, callerSkip ...int) {
		if message != "" {
			t.Error(message)
		}
	}
}
