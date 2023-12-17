package main

import (
	"io/ioutil"
	"strings"
)

type deckFileStorage struct{}

func (dp deckFileStorage) readFromFile(fileName string) (deck, error) {
	fileContentBytes, err := ioutil.ReadFile(fileName)
	var d deck
	if err == nil {
		fileContent := strings.Split(string(fileContentBytes), ",")
		var cards []card
		for _, cardString := range fileContent {
			cards = append(cards, card(cardString))
		}
		d = deck(cards)
	}

	return d, err
}

func (dp deckFileStorage) writeToFile(fileName string, d deck) error {
	fileContent := ""
	for idx, card := range d {
		if idx == len(d)-1 {
			fileContent += card.toString()
		} else {
			fileContent += card.toString() + ","
		}
	}
	// write the whole body at once
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	return err
}
