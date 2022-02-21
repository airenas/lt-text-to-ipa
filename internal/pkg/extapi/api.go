package extapi

//AccentOutputElement accentuate method output element
type AccentOutputElement struct {
	Accent []Accent `json:"accent"`
	Word   string   `json:"word"`
	Error  string   `json:"error,omitempty"`
}

//Accent structure for one word accent
type Accent struct {
	MF       string          `json:"mf"`
	Mi       string          `json:"mi"`
	MiVdu    string          `json:"mi_vdu"`
	Mih      string          `json:"mih"`
	Error    string          `json:"error,omitempty"`
	Variants []AccentVariant `json:"variants"`
}

//AccentVariant - accenters's result
type AccentVariant struct {
	Accent   int     `json:"accent"`
	Accented string  `json:"accented,omitempty"`
	Ml       string  `json:"ml"`
	Syll     string  `json:"syll"`
	Usage    float64 `json:"usage,omitempty"`
	Meaning  string  `json:"meaning,omitempty"`
}

//TransInput transcription input
type TransInput struct {
	Word string `json:"word"`
	Syll string `json:"syll"`
	User string `json:"user"`
	Ml   string `json:"ml"`
	Rc   string `json:"rc"`
	Acc  int    `json:"acc"`
}

//TransOutput transcription output
type TransOutput struct {
	Transcription []Trans `json:"transcription"`
	Word          string  `json:"word"`
	Error         string  `json:"error"`
}

//Trans transcription variant
type Trans struct {
	Transcription string `json:"transcription"`
}

//IPAInput input for IPA method
type IPAInput struct {
	Transcription string `json:"transcription"`
}

//IPAOutput output of IPA method
type IPAOutput struct {
	Transcription string `json:"transcription"`
	IPA           string `json:"ipa"`
}
