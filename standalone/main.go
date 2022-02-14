package main

import (
  "flag"
  "log"
  "fmt"
  "os"

  m2n "github.com/Th3-S1lenc3/morse2note"
)

func main() {
  morse2note, err := m2n.NewMorse2Note()
  if err != nil {
    log.Fatal(err)
  }

  morseString := flag.String("m", "", "Morse code string to convert.")
  appDir := flag.String("appDir", "", "Path to app config dir, if part of larger application.")
  outputDir := flag.String("o", "", "Path to write output file.")
  override := flag.Bool("override", false, "Override existing JSON file.")
  indent := flag.Bool("indent", false, "Indent output JSON.")

  flag.Parse()

  err = morse2note.Init(*morseString, 4, *appDir)
  if err != nil {
    log.Fatal(err)
  }

  err = morse2note.Convert()
  if err != nil {
    log.Fatal(err)
  }

  var outputPath string

  if *outputDir == "" {
    cwd, err := os.Getwd()
    if err != nil {
      log.Fatal(err)
    }

    outputPath = fmt.Sprintf("%s/note.json", cwd)
  } else {
    outputPath = fmt.Sprintf("%s/note.json", *outputDir)
  }

  err = morse2note.WriteNotesToFile(outputPath, *override, *indent)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Done")
}
