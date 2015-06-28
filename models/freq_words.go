package models

type FreqWords struct {
  Words map[string]int
}

func (fw *FreqWords) AddWord(word string) {
  value, _ := fw.Words[word]
  fw.Words[word] = value + 1
}