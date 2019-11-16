package main

import (
	"fmt"
	"github.com/keewonma/deck"
	"strings"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			//Ace is worth 1 or 11. currently 1. now max = 11
			return minScore + 10
		}
	}
	return minScore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func main() {
	cards := deck.New(deck.Deck(3), deck.Shuffle)
	var card deck.Card
	var player, dealer Hand
	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} {
			card, cards = draw(cards)
			*hand = append(*hand, card) //pointer to hand

		}
	}
	var input string
	for input != "s" {
		fmt.Printf("***CURRENT HAND***\n")
		fmt.Printf("\n")
		fmt.Println("Player:", player, "\nScore:", player.Score())
		fmt.Printf("\n")
		fmt.Println("Dealer:", dealer.DealerString())
		fmt.Println("What do you want to do? (h)it, (s)tand")
		fmt.Scanf("%s/n", &input)
		switch input {
		case "h":
			card, cards = draw(cards)
			player = append(player, card)
		}
	}
	//If dealer score <= 16, dealer hits
	//If dealer has a soft 17, dealer hits
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
		card, cards = draw(cards)
		dealer = append(dealer, card)

	}

	pScore, dScore := player.Score(), dealer.Score()
	fmt.Printf("---FINAL HANDS---\n")
	fmt.Printf("\n")
	fmt.Println("Player: ", player, "\nScore:", pScore)
	fmt.Printf("\n")
	fmt.Println("Dealer: ", dealer, "\nScore:", dScore)
	switch {
	case pScore > 21:
		//Player Over 21
		fmt.Println("You busted")
	case dScore > 21:
		//Dealer Over 21
		fmt.Println("Dealer busted")
	case pScore > dScore:
		//Player wins by having a better score than the dealer
		fmt.Println("You win!")
	case dScore == pScore:
		//Draw
		fmt.Println("Push")
	case dScore > pScore:
		//Dealer wins by having a better score than the dealer
		fmt.Println("Push")
	}

}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
