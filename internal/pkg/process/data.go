package process

import "github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"

// Data working data for one request
type Data struct {
	OriginalText string
	Text         string

	Words []*ProcessedWord

	Result []*api.ResultWord
}

//ProcessedWord keeps one word info
type ProcessedWord struct {
	Tagged             TaggedWord
	AccentVariant      *AccentVariant
	AccentCount        int
	TranscriptionCount int
	Transcription      string
	IPA                string
	Clitic             *Clitic
}

//StringTypeEnum represent possible string types
type StringTypeEnum int

const (
	//Word value
	Word StringTypeEnum = iota + 1
	//OtherWord value
	OtherWord
	//Separator value
	Separator
	//Space value
	Space
	//SentenceEnd value - data normalized by user
	SentenceEnd
)

//TaggedWord - tagger's result
type TaggedWord struct {
	Type   StringTypeEnum
	String string
	Mi     string
	Lemma  string
}

//AccentVariant - accenters's result
type AccentVariant struct {
	Accent   int     `json:"accent"`
	Accented string  `json:"accented"`
	Ml       string  `json:"ml"`
	Syll     string  `json:"syll"`
	Usage    float64 `json:"usage"`
}

//Clitic contains clitic analysis data
type Clitic struct {
	Accent       int
	AccentedType string
	Type         string
	Pos          int
}
