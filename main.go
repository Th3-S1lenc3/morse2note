package main

import (
  "flag"
  "log"
)

func main() {
  convert, err := NewConvert()
  if err != nil {
    log.Fatal(err)
  }

  morseString := flag.String("m", "", "Morse code string to convert.")

  flag.Parse()

  err = convert.Init(*morseString, 4)
  if err != nil {
    log.Fatal(err)
  }

  err = convert.Convert2Note()
  if err != nil {
    log.Fatal(err)
  }

  err = convert.WriteNotesToFile("", true, true)
  if err != nil {
    log.Fatal(err)
  }
}
