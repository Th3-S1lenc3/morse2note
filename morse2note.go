package morse2note

import (
  "strings"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "math"
  "strconv"
  "os"
  "github.com/cavaliergopher/grab/v3"
  "time"
)

type Morse2Note struct {
  morseString string
  startingOctave int
  dictionary Notes
  notes CNotes
}

func NewMorse2Note() *Morse2Note {
  return &Morse2Note{}
}

func (m *Morse2Note) Convert() error {
  morseWholeWords := strings.Split(m.morseString, "//")

  for i := 0; i < len(morseWholeWords); i++ {
    morseWords := strings.Split(morseWholeWords[i], "/")

    for j := 0; j < len(morseWords); j++ {
      morseChars := strings.Split(morseWords[j], "")

      block := CBlock{}

      for k := 0; k < len(morseChars); k++ {
        char := morseChars[k]

        octave := m.startingOctave
        noteType := "quarterBeat"

        startingIndex := -2 - k

        if char == "-" {
          octave -= 1
          startingIndex = int(math.Abs(float64(startingIndex)))
        }

        note, err := m.findNoteInDictionary(startingIndex, octave)
        if err != nil {
          return err
        }

        note.NoteType = noteType

        block.Notes = append(block.Notes, note)
      }

      m.notes.Blocks = append(m.notes.Blocks, block)
    }
  }


  return nil
}

func (m *Morse2Note) GetDictionary() (Notes, error) {
  return m.dictionary, nil
}

func (m *Morse2Note) GetNotes() (CNotes, error) {
  return m.notes, nil
}

func (m *Morse2Note) WriteNotesToFile(filePath string, indent bool, override bool) error {
  if filePath == "" {
    filePath = "note.json"
  }

  if indent == true {
    data, err := json.MarshalIndent(m.notes, "", "  ")
    if err != nil {
      return err
    }

    return safeWrite(filePath, data, os.FileMode(0644), override)
  } else {
    data, err := json.Marshal(m.notes)
    if err != nil {
      return err
    }

    return safeWrite(filePath, data, os.FileMode(0644), override)
  }

  return nil
}

func (m *Morse2Note) findNoteInDictionary(index int, octave int) (Note, error) {
  notes := m.dictionary.Notes
  piano := m.dictionary.Piano

  middleCIndex := len(notes) / 2

  noteIndex := middleCIndex + index

  noteLetter := notes[noteIndex]

  for i := 0; i < len(piano); i++ {
    note := piano[i]

    isCorrectNote := note.Note == noteLetter
    isCorrectOctave := note.Octave == strconv.Itoa(octave)

    if isCorrectNote && isCorrectOctave {
      return note, nil
    }
  }

  return Note{}, fmt.Errorf("Not Found.")
}

func (m *Morse2Note) checkMorseString() error {
  morseWholeWords := strings.Split(m.morseString, "//")

  for i := 0; i < len(morseWholeWords); i++ {
    morseWords := strings.Split(morseWholeWords[i], "/")

    for j := 0; j < len(morseWords); j++ {
      morseChars := strings.Split(morseWords[j], "")

      for k := 0; k < len(morseChars); k++ {
        if morseChars[i] != "." && morseChars[i] != "-" {
          return fmt.Errorf("Invalid Morse String: %s; %s; %s", m.morseString, morseChars, morseChars[i])
        }
      }
    }
  }

  return nil
}

func (m *Morse2Note) DownloadNotes(configDir string, fileName string) error {
  fmt.Printf("Cannot find \"%s\" in \"%s\"\n", fileName, configDir)

  remoteFileURL := "https://raw.githubusercontent.com/Th3-S1lenc3/morse2note/master/json/notes.min.json"

  // Create Client
  client := grab.NewClient()
  req, _ := grab.NewRequest(configDir, remoteFileURL)

  // Start Download
  fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

  // Start UI Loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
  for {
    select {
    case <-t.C:
      fmt.Printf(
        "  transferred %v / %v bytes (%.2f%%)\n",
        resp.BytesComplete(),
        resp.Size,
        100 * resp.Progress(),
      )
    case <- resp.Done:
      break Loop
    }
  }

  if err := resp.Err(); err != nil {
    return fmt.Errorf("Download failed %v\n", err)
  }

  fmt.Printf("Download saved to %v \n", resp.Filename)

  return nil
}

func (m *Morse2Note) Init(morseString string, startingOctave int, appDir string) error {
  m.morseString = strings.ReplaceAll(morseString, " ", "")
  m.morseString = strings.Trim(m.morseString, "/")

  err := m.checkMorseString()
  if err != nil {
    return err
  }

  m.startingOctave = startingOctave

  m.notes = CNotes{}

  if appDir == "" {
    appDir, err = os.UserConfigDir()
    if err != nil {
      return err
    }
  }

  _, err = os.Stat(appDir)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Cannot find directory: \"%s\"", appDir)
	}

  configDir := fmt.Sprintf("%s/Morse2Note", appDir)

  _, err = os.Stat(configDir)
  if err != nil && os.IsNotExist(err) {
    err = os.Mkdir(configDir, os.FileMode(0755))
    if err != nil {
      return err
    }
	}

  notesJsonFilePath := fmt.Sprintf("%s/notes.min.json", configDir)

  _, err = os.Stat(notesJsonFilePath)
  if err != nil && os.IsNotExist(err) {
    err := m.DownloadNotes(configDir, "notes.min.json")
    if err != nil {
      return err
    }
  }

  jsonData, err := ioutil.ReadFile(notesJsonFilePath)
  if err != nil {
    return err
  }

  err = json.Unmarshal(jsonData, &m.dictionary)
  if err != nil {
    return err
  }

  return nil
}
