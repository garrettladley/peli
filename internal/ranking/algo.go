package ranking

import (
	"math"

	"github.com/garrettladley/peli/internal/db"
)

const K = 32.0 // Elo K-factor

// CompareResult is returned by UI after user makes a choice.
type CompareResult struct {
	WinnerID   int64
	WinnerName string
	LoserID    int64
	LoserName  string
}

// BinarySearchState tracks the current state of binary insertion.
type BinarySearchState struct {
	NewRestaurant db.Restaurant
	Ranked        []db.Restaurant // sorted best-to-worst
	Low           int
	High          int
	CurrentMid    int
	Comparisons   []CompareResult
	Done          bool
	InsertAt      int // final position (0-indexed)
}

// NewBinarySearch creates a new binary search state for inserting a restaurant.
func NewBinarySearch(newR db.Restaurant, ranked []db.Restaurant) *BinarySearchState {
	if len(ranked) == 0 {
		return &BinarySearchState{
			NewRestaurant: newR,
			Done:          true,
			InsertAt:      0,
		}
	}
	bs := &BinarySearchState{
		NewRestaurant: newR,
		Ranked:        ranked,
		Low:           0,
		High:          len(ranked) - 1,
	}
	bs.CurrentMid = (bs.Low + bs.High) / 2
	return bs
}

// Step processes user's choice and advances search.
// newWins is true if the new restaurant is preferred over the current opponent.
func (bs *BinarySearchState) Step(newWins bool) {
	opponent := bs.Ranked[bs.CurrentMid]

	if newWins {
		bs.Comparisons = append(bs.Comparisons, CompareResult{
			WinnerID:   bs.NewRestaurant.ID,
			WinnerName: bs.NewRestaurant.Name,
			LoserID:    opponent.ID,
			LoserName:  opponent.Name,
		})
		bs.High = bs.CurrentMid - 1
	} else {
		bs.Comparisons = append(bs.Comparisons, CompareResult{
			WinnerID:   opponent.ID,
			WinnerName: opponent.Name,
			LoserID:    bs.NewRestaurant.ID,
			LoserName:  bs.NewRestaurant.Name,
		})
		bs.Low = bs.CurrentMid + 1
	}

	if bs.Low > bs.High {
		bs.Done = true
		bs.InsertAt = bs.Low // 0-indexed
	} else {
		bs.CurrentMid = (bs.Low + bs.High) / 2
	}
}

// CurrentOpponent returns the restaurant to compare against.
func (bs *BinarySearchState) CurrentOpponent() db.Restaurant {
	return bs.Ranked[bs.CurrentMid]
}

// ExpectedComparisons returns ceil(log2(n+1)).
func (bs *BinarySearchState) ExpectedComparisons() int {
	n := len(bs.Ranked) + 1
	count := 0
	for n > 0 {
		count++
		n >>= 1
	}
	return count
}

// ComparisonNumber returns the current comparison number (1-indexed).
func (bs *BinarySearchState) ComparisonNumber() int {
	return len(bs.Comparisons) + 1
}

// UpdateElo updates the Elo ratings for winner and loser.
func UpdateElo(winner, loser *db.Restaurant) {
	expectedWinner := 1.0 / (1.0 + math.Pow(10, (loser.Elo-winner.Elo)/400))
	expectedLoser := 1.0 - expectedWinner

	winner.Elo += K * (1.0 - expectedWinner)
	loser.Elo += K * (0.0 - expectedLoser)

	// normalize to 0-10 display score (sigmoid)
	winner.Score = 10.0 / (1.0 + math.Exp(-(winner.Elo-1500)/200))
	loser.Score = 10.0 / (1.0 + math.Exp(-(loser.Elo-1500)/200))
}

// CalculateScore computes the 0-10 display score from Elo.
func CalculateScore(elo float64) float64 {
	return 10.0 / (1.0 + math.Exp(-(elo-1500)/200))
}

// RecalculateAllScores recalculates scores for all restaurants based on position.
// This is an alternative scoring method based purely on rank position.
func RecalculateAllScores(restaurants []db.Restaurant) {
	n := len(restaurants)
	if n == 0 {
		return
	}
	for i := range restaurants {
		// position from bottom (n for #1, 1 for last)
		posFromBottom := float64(n - i)
		restaurants[i].Score = 10.0 * (posFromBottom / float64(n))
	}
}
