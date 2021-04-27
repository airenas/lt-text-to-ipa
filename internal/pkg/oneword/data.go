package oneword

import "github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"

// Data working data for one request
type Data struct {
	Word  string
	Words []*WorkingWord

	Result *api.WordInfo
}

type WorkingWord struct {
	Meaning string
	MI      string
	Accent  int
	Syll    string
	Lemma   string

	Transcriptions []string
	IPAs           []string
}
