package main

func main() {
	d := newDeck()
	d.shuffle()
	// d.print()

	var dp deckStorage = deckFileStorage{}
	dp.writeToFile("deck", d)
	d, _ = dp.readFromFile("deck")
	d.print()
}
