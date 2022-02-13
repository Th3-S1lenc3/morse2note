package main

import (
  "flag"
  "log"
  "fmt"
)

func main() {
  convert, err := NewConvert()
  if err != nil {
    log.Fatal(err)
  }

  morseString := flag.String("m", "", "Morse code string to convert.")
  appDir := flag.String("appDir", "", "Path to app config dir, if part of larger application")

  flag.Parse()

  err = convert.Init(*morseString, 4, *appDir)
  if err != nil {
    log.Fatal(err)
  }

  err = convert.Convert()
  if err != nil {
    log.Fatal(err)
  }

  err = convert.WriteNotesToFile("", true, true)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Done")
}
