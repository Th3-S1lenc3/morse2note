package main

import (
  "strings"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "math"
  "strconv"
  "os"
)

type Convert struct {
  morseString string
  startingOctave int
  dictionary Notes
  notes CNotes
}

func NewConvert() (*Convert, error) {
  return &Convert{}, nil
}

func (c *Convert) Convert() error {
  morseWholeWords := strings.Split(c.morseString, "//")

  for i := 0; i < len(morseWholeWords); i++ {
    morseWords := strings.Split(morseWholeWords[i], "/")

    for j := 0; j < len(morseWords); j++ {
      morseChars := strings.Split(morseWords[j], "")

      block := CBlock{}

      for k := 0; k < len(morseChars); k++ {
        char := morseChars[k]

        octave := c.startingOctave
        noteType := "quarterBeat"

        startingIndex := -2 - k

        if char == "-" {
          octave -= 1
          startingIndex = int(math.Abs(float64(startingIndex)))
        }

        note, err := c.findNoteInDictionary(startingIndex, octave)
        if err != nil {
          return err
        }

        note.NoteType = noteType

        block.Notes = append(block.Notes, note)
      }

      c.notes.Blocks = append(c.notes.Blocks, block)
    }
  }


  return nil
}

func (c *Convert) GetDictionary() (Notes, error) {
  return c.dictionary, nil
}

func (c *Convert) GetNotes() (CNotes, error) {
  return c.notes, nil
}

func (c *Convert) WriteNotesToFile(filePath string, indent bool, override bool) error {
  if filePath == "" {
    filePath = "note.json"
  }

  if indent == true {
    data, err := json.MarshalIndent(c.notes, "", "  ")
    if err != nil {
      return err
    }

    return safeWrite(filePath, data, os.FileMode(0644), override)
  } else {
    data, err := json.Marshal(c.notes)
    if err != nil {
      return err
    }

    return safeWrite(filePath, data, os.FileMode(0644), override)
  }

  return nil
}

func (c *Convert) findNoteInDictionary(index int, octave int) (Note, error) {
  notes := c.dictionary.Notes
  piano := c.dictionary.Piano

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

func (c *Convert) checkMorseString() error {
  morseWholeWords := strings.Split(c.morseString, "//")

  for i := 0; i < len(morseWholeWords); i++ {
    morseWords := strings.Split(morseWholeWords[i], "/")

    for j := 0; j < len(morseWords); j++ {
      morseChars := strings.Split(morseWords[j], "")

      for k := 0; k < len(morseChars); k++ {
        if morseChars[i] != "." && morseChars[i] != "-" {
          return fmt.Errorf("Invalid Morse String: %s; %s; %s", c.morseString, morseChars, morseChars[i])
        }
      }
    }
  }

  return nil
}

func (c *Convert) Init(morseString string, startingOctave int) error {
  c.morseString = strings.ReplaceAll(morseString, " ", "")
  c.morseString = strings.Trim(c.morseString, "/")

  err := c.checkMorseString()
  if err != nil {
    return err
  }

  c.startingOctave = startingOctave

  c.notes = CNotes{}

  jsonData, err := ioutil.ReadFile("notes.json")
  if err != nil {
    return err
  }

  err = json.Unmarshal(jsonData, &c.dictionary)
  if err != nil {
    return err
  }

  return nil
}
