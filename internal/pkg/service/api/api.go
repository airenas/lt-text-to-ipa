package api

// ResultWord is service output
type ResultWord struct {
	Type    string    `json:"type"`
	String  string    `json:"string,omitempty"`
	IPA     string    `json:"ipa,omitempty"`
	IPAType string    `json:"ipaType,omitempty"`
	Info    *WordInfo `json:"info,omitempty"`
}

// WordInfo is a service output for one word
type WordInfo struct {
	Word           string          `json:"word"`
	Transcriptions []Transcription `json:"transcription"`
}

// Transcription is a service output for one transcription variant
type Transcription struct {
	IPAs        []string `json:"ipas"`
	Information []MIInfo `json:"information"`
}

// MIInfo is a keeper for mi and meaning of a word
type MIInfo struct {
	MI      string `json:"mi"`
	Meaning string `json:"meaning,omitempty"`
}
