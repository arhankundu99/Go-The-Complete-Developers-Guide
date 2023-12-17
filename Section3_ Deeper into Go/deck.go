package main

import (
	"math/rand"
)

type deck []card

// Try to follow this convention when you are initialising a struct using a function
// the function name should be of this format: newStruct
func newDeck() deck {
	cardValues := [12]string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Jack", "Queen", "King"}
	cardTypes := [4]string{"Spades", "Clubs", "Hearts", "Diamonds"}

	nd := make([]card, 0)
	for _, cardValue := range cardValues {
		for _, cardType := range cardTypes {
			nd = append(nd, card(cardValue+" Of "+cardType))
		}
	}
	return nd
}

func (d deck) print() {
	for _, card := range d {
		card.print()
	}
}

func (d deck) shuffle() {
	for i := 0; i < len(d); i++ {
		randomIdx := rand.Intn(len(d))
		d[i], d[randomIdx] = d[randomIdx], d[i]
	}
}
