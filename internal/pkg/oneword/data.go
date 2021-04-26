package oneword

import "github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"

// Data working data for one request
type Data struct {
	Word string

	Result *api.WordInfo
}
