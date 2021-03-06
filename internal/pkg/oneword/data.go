package oneword

import "github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"

// Data working data for one word processing
type Data struct {
	Word        string
	ReturnSAMPA bool

	Words []*WorkingWord

	Result *api.WordInfo
}

//WorkingWord structure for one working word
type WorkingWord struct {
	Meaning string
	MI      string
	Accent  int
	Syll    string
	Lemma   string

	Transcriptions []string
	IPAs           []string
}
