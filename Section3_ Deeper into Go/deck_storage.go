package main

type deckStorage interface {
	readFromFile(fileName string) (deck, error)
	writeToFile(fileName string, d deck) error
}
