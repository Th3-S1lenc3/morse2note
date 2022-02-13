package morse2note

type Notes struct {
  Notes []string `json:"notes"`
  Piano []Note `json:"piano"`
}

type Note struct {
  Note string `json:"note"`
  Octave string `json:"octave"`
  Frequency Frequency `json:"frequency"`
  NoteType string `json:"type"`
}

type Frequency struct {
  Natural string `json:"natural"`
  Flat string `json:"flat"`
  Sharp string `json:"sharp"`
}

type CNotes struct {
  Blocks []CBlock `json:"blocks"`
}

type CBlock struct {
  Notes []Note `json:"notes"`
}
