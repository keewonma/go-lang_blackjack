package main

import (
	"fmt"
	"github.com/keewonma/deck"
	"strings"
)

//Hand slice
type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

//DealerString - Hides Dealers Second Card
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

//Score - Calculates Score
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

//MinScore - Calculate minimum score with ace = 1 and kings, queens, jacks = 10
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Shuffle - Shuffle the Deck
func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

//Deal - Deal 2 cards, with a max of 5 cards in hand
func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

//Hit - Deals a card to hand
func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

//Stand - moves to next state.
func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++ //move to dealer turn
	return ret
}

//EndHand - ends the hand and prints final score
func EndHand(gs GameState) GameState {
	ret := clone(gs)
	//Calculate Score
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Printf("---FINAL HANDS---\n")
	fmt.Printf("\n")
	fmt.Println("Player: ", ret.Player, "\nScore:", pScore)
	fmt.Printf("\n")
	fmt.Println("Dealer: ", ret.Dealer, "\nScore:", dScore)
	fmt.Println()
	switch {

	case pScore > 21:
		//Player Over 21
		fmt.Println("-----You busted!!!-----")
	case dScore > 21:
		//Dealer Over 21
		fmt.Println("-----Dealer busted!!!-----")
	case pScore > dScore:
		//Player wins by having a better score than the dealer
		fmt.Println("-----You win!!!-----")
	case dScore > pScore:
		//Dealer wins by having a better score than the dealer
		fmt.Println("-----You lose!!!-----")
	case dScore == pScore:
		//Draw
		fmt.Println("-----Push-----")
	}
	fmt.Println("\n---NEW GAME---")
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {
	//initial starting dealing pile with 3 decks and random shuffle
	var gs GameState
	gs = Shuffle(gs)

	for i := 0; i < 10; i++ {
		//Deal two cards
		gs = Deal(gs)
		var input string
		for gs.State == StatePlayerTurn {
			pScore := gs.Player.Score()
			fmt.Println("Player: ", gs.Player, "\nScore:", pScore)
			fmt.Println("Dealer:", gs.Dealer.DealerString())
			fmt.Println("What will you do? (h)it, (s)tand")
			fmt.Printf("\n")
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid option:", input)
			}
		}

		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

//State - gamestate
type State int8

const (
	//StatePlayerTurn - player's turn
	StatePlayerTurn State = iota
	//StateDealerTurn - dealer's turn
	StateDealerTurn
	//StateHandOver - hand is over
	StateHandOver
)

//GameState - struct of gamestates
type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

//CurrentPlayer - whose turn it is
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}
