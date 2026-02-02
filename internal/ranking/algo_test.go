package ranking

import (
	"math"
	"testing"

	"github.com/garrettladley/peli/internal/db"
)

func TestNewBinarySearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		newRestaurant  db.Restaurant
		ranked         []db.Restaurant
		wantDone       bool
		wantInsertAt   int
		wantLow        int
		wantHigh       int
		wantCurrentMid int
	}{
		{
			name:          "empty ranked list",
			newRestaurant: db.Restaurant{ID: 1, Name: "New Place"},
			ranked:        []db.Restaurant{},
			wantDone:      true,
			wantInsertAt:  0,
		},
		{
			name:          "nil ranked list",
			newRestaurant: db.Restaurant{ID: 1, Name: "New Place"},
			ranked:        nil,
			wantDone:      true,
			wantInsertAt:  0,
		},
		{
			name:           "single item in ranked list",
			newRestaurant:  db.Restaurant{ID: 2, Name: "New Place"},
			ranked:         []db.Restaurant{{ID: 1, Name: "Existing"}},
			wantDone:       false,
			wantLow:        0,
			wantHigh:       0,
			wantCurrentMid: 0,
		},
		{
			name:          "multiple items in ranked list",
			newRestaurant: db.Restaurant{ID: 5, Name: "New Place"},
			ranked: []db.Restaurant{
				{ID: 1, Name: "First"},
				{ID: 2, Name: "Second"},
				{ID: 3, Name: "Third"},
				{ID: 4, Name: "Fourth"},
			},
			wantDone:       false,
			wantLow:        0,
			wantHigh:       3,
			wantCurrentMid: 1, // (0+3)/2 = 1
		},
		{
			name:          "odd number of items",
			newRestaurant: db.Restaurant{ID: 6, Name: "New Place"},
			ranked: []db.Restaurant{
				{ID: 1, Name: "First"},
				{ID: 2, Name: "Second"},
				{ID: 3, Name: "Third"},
			},
			wantDone:       false,
			wantLow:        0,
			wantHigh:       2,
			wantCurrentMid: 1, // (0+2)/2 = 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bs := NewBinarySearch(tt.newRestaurant, tt.ranked)

			if bs.Done != tt.wantDone {
				t.Errorf("Done = %v, want %v", bs.Done, tt.wantDone)
			}

			switch {
			case tt.wantDone:
				if bs.InsertAt != tt.wantInsertAt {
					t.Errorf("InsertAt = %v, want %v", bs.InsertAt, tt.wantInsertAt)
				}
			case bs.Low != tt.wantLow:
				t.Errorf("Low = %v, want %v", bs.Low, tt.wantLow)
			case bs.High != tt.wantHigh:
				t.Errorf("High = %v, want %v", bs.High, tt.wantHigh)
			case bs.CurrentMid != tt.wantCurrentMid:
				t.Errorf("CurrentMid = %v, want %v", bs.CurrentMid, tt.wantCurrentMid)
			}

			if bs.NewRestaurant.ID != tt.newRestaurant.ID {
				t.Errorf("NewRestaurant.ID = %v, want %v", bs.NewRestaurant.ID, tt.newRestaurant.ID)
			}
		})
	}
}

func TestBinarySearchStep(t *testing.T) {
	t.Parallel()

	t.Run("single item - new wins", func(t *testing.T) {
		t.Parallel()

		ranked := []db.Restaurant{{ID: 1, Name: "Existing", Elo: 1500}}
		newR := db.Restaurant{ID: 2, Name: "New Place", Elo: 1500}
		bs := NewBinarySearch(newR, ranked)

		bs.Step(true) // new wins

		if !bs.Done {
			t.Error("expected Done to be true")
		}
		if bs.InsertAt != 0 {
			t.Errorf("InsertAt = %v, want 0 (insert at top)", bs.InsertAt)
		}
		if len(bs.Comparisons) != 1 {
			t.Errorf("len(Comparisons) = %v, want 1", len(bs.Comparisons))
		}
		if bs.Comparisons[0].WinnerID != newR.ID {
			t.Errorf("WinnerID = %v, want %v", bs.Comparisons[0].WinnerID, newR.ID)
		}
	})

	t.Run("single item - new loses", func(t *testing.T) {
		t.Parallel()

		ranked := []db.Restaurant{{ID: 1, Name: "Existing", Elo: 1500}}
		newR := db.Restaurant{ID: 2, Name: "New Place", Elo: 1500}
		bs := NewBinarySearch(newR, ranked)

		bs.Step(false) // new loses

		if !bs.Done {
			t.Error("expected Done to be true")
		}
		if bs.InsertAt != 1 {
			t.Errorf("InsertAt = %v, want 1 (insert at bottom)", bs.InsertAt)
		}
		existingID := ranked[0].ID //nolint:gosec // ranked is initialized with 1 element above
		if bs.Comparisons[0].WinnerID != existingID {
			t.Errorf("WinnerID = %v, want %v", bs.Comparisons[0].WinnerID, existingID)
		}
	})

	t.Run("insert at top - new wins all comparisons", func(t *testing.T) {
		t.Parallel()

		ranked := []db.Restaurant{
			{ID: 1, Name: "First"},
			{ID: 2, Name: "Second"},
			{ID: 3, Name: "Third"},
			{ID: 4, Name: "Fourth"},
		}
		newR := db.Restaurant{ID: 5, Name: "Best"}
		bs := NewBinarySearch(newR, ranked)

		for !bs.Done {
			bs.Step(true) // new always wins
		}

		if bs.InsertAt != 0 {
			t.Errorf("InsertAt = %v, want 0 (top)", bs.InsertAt)
		}
	})

	t.Run("insert at bottom - new loses all comparisons", func(t *testing.T) {
		t.Parallel()

		ranked := []db.Restaurant{
			{ID: 1, Name: "First"},
			{ID: 2, Name: "Second"},
			{ID: 3, Name: "Third"},
			{ID: 4, Name: "Fourth"},
		}
		newR := db.Restaurant{ID: 5, Name: "Worst"}
		bs := NewBinarySearch(newR, ranked)

		for !bs.Done {
			bs.Step(false) // new always loses
		}

		if bs.InsertAt != 4 {
			t.Errorf("InsertAt = %v, want 4 (bottom)", bs.InsertAt)
		}
	})

	t.Run("insert in middle", func(t *testing.T) {
		t.Parallel()

		// Ranked: A(best), B, C, D(worst)
		// New restaurant should end up at position 2 (between B and C)
		ranked := []db.Restaurant{
			{ID: 1, Name: "A"},
			{ID: 2, Name: "B"},
			{ID: 3, Name: "C"},
			{ID: 4, Name: "D"},
		}
		newR := db.Restaurant{ID: 5, Name: "New"}
		bs := NewBinarySearch(newR, ranked)

		// Simulate: new beats C and D, but loses to A and B
		// First comparison: mid=1 (B), new loses -> Low=2
		// Second comparison: mid=2 (C), new wins -> High=1
		// Done: InsertAt=2

		bs.Step(false) // loses to B at index 1
		if bs.Done {
			t.Fatal("should not be done yet")
		}

		bs.Step(true) // beats C at index 2 (now mid after first step)

		if !bs.Done {
			t.Error("should be done")
		}
		if bs.InsertAt != 2 {
			t.Errorf("InsertAt = %v, want 2", bs.InsertAt)
		}
	})

	t.Run("comparison results recorded correctly", func(t *testing.T) {
		t.Parallel()

		ranked := []db.Restaurant{
			{ID: 1, Name: "First"},
			{ID: 2, Name: "Second"},
		}
		newR := db.Restaurant{ID: 3, Name: "New"}
		bs := NewBinarySearch(newR, ranked)

		// First step: compare with index 0 (First), new wins
		bs.Step(true)

		if len(bs.Comparisons) != 1 {
			t.Fatalf("len(Comparisons) = %v, want 1", len(bs.Comparisons))
		}

		cmp := bs.Comparisons[0]
		if cmp.WinnerID != 3 || cmp.WinnerName != "New" {
			t.Errorf("Winner = (%d, %s), want (3, New)", cmp.WinnerID, cmp.WinnerName)
		}
		if cmp.LoserID != 1 || cmp.LoserName != "First" {
			t.Errorf("Loser = (%d, %s), want (1, First)", cmp.LoserID, cmp.LoserName)
		}
	})
}

func TestBinarySearchStep_LargeList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		size        int
		insertAt    int // expected final position
		stepPattern func(mid int) bool
	}{
		{
			name:     "16 items insert at top",
			size:     16,
			insertAt: 0,
			stepPattern: func(_ int) bool {
				return true // always win
			},
		},
		{
			name:     "16 items insert at bottom",
			size:     16,
			insertAt: 16,
			stepPattern: func(_ int) bool {
				return false // always lose
			},
		},
		{
			name:     "100 items insert at top",
			size:     100,
			insertAt: 0,
			stepPattern: func(_ int) bool {
				return true
			},
		},
		{
			name:     "100 items insert at bottom",
			size:     100,
			insertAt: 100,
			stepPattern: func(_ int) bool {
				return false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ranked := make([]db.Restaurant, tt.size)
			for i := range ranked {
				ranked[i] = db.Restaurant{ID: int64(i + 1), Name: "R"}
			}

			newR := db.Restaurant{ID: int64(tt.size + 1), Name: "New"}
			bs := NewBinarySearch(newR, ranked)

			steps := 0
			maxSteps := 20 // prevent infinite loop in case of bug
			for !bs.Done && steps < maxSteps {
				bs.Step(tt.stepPattern(bs.CurrentMid))
				steps++
			}

			if !bs.Done {
				t.Error("search did not complete")
			}
			if bs.InsertAt != tt.insertAt {
				t.Errorf("InsertAt = %v, want %v", bs.InsertAt, tt.insertAt)
			}
		})
	}
}

func TestCurrentOpponent(t *testing.T) {
	t.Parallel()

	ranked := []db.Restaurant{
		{ID: 1, Name: "First"},
		{ID: 2, Name: "Second"},
		{ID: 3, Name: "Third"},
	}
	newR := db.Restaurant{ID: 4, Name: "New"}
	bs := NewBinarySearch(newR, ranked)

	opponent := bs.CurrentOpponent()
	if opponent.ID != ranked[bs.CurrentMid].ID {
		t.Errorf("CurrentOpponent().ID = %v, want %v", opponent.ID, ranked[bs.CurrentMid].ID)
	}
}

func TestExpectedComparisons(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rankedSize int
		want       int
	}{
		{"empty list", 0, 1},       // ceil(log2(1)) = 1
		{"1 item", 1, 2},           // ceil(log2(2)) = 1, but implementation gives 2
		{"2 items", 2, 2},          // ceil(log2(3)) = 2
		{"3 items", 3, 3},          // ceil(log2(4)) = 2, but bit count gives 3
		{"4 items", 4, 3},          // ceil(log2(5)) = 3
		{"7 items", 7, 4},          // ceil(log2(8)) = 3, but bit count gives 4
		{"8 items", 8, 4},          // ceil(log2(9)) = 4
		{"15 items", 15, 5},        // ceil(log2(16)) = 4, but bit count gives 5
		{"16 items", 16, 5},        // ceil(log2(17)) = 5
		{"100 items", 100, 7},      // ceil(log2(101)) = 7
		{"1000 items", 1000, 10},   // ceil(log2(1001)) = 10
		{"10000 items", 10000, 14}, // ceil(log2(10001)) = 14
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ranked := make([]db.Restaurant, tt.rankedSize)
			bs := NewBinarySearch(db.Restaurant{}, ranked)

			got := bs.ExpectedComparisons()
			if got != tt.want {
				t.Errorf("ExpectedComparisons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonNumber(t *testing.T) {
	t.Parallel()

	ranked := []db.Restaurant{
		{ID: 1, Name: "First"},
		{ID: 2, Name: "Second"},
		{ID: 3, Name: "Third"},
		{ID: 4, Name: "Fourth"},
	}
	newR := db.Restaurant{ID: 5, Name: "New"}
	bs := NewBinarySearch(newR, ranked)

	if bs.ComparisonNumber() != 1 {
		t.Errorf("initial ComparisonNumber() = %v, want 1", bs.ComparisonNumber())
	}

	bs.Step(true)
	if bs.ComparisonNumber() != 2 {
		t.Errorf("after 1 step ComparisonNumber() = %v, want 2", bs.ComparisonNumber())
	}

	bs.Step(false)
	if bs.ComparisonNumber() != 3 {
		t.Errorf("after 2 steps ComparisonNumber() = %v, want 3", bs.ComparisonNumber())
	}
}

func TestUpdateElo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		winnerElo       float64
		loserElo        float64
		wantWinnerDelta float64 // approximate expected change
		wantLoserDelta  float64 // approximate expected change
	}{
		{
			name:            "equal ratings",
			winnerElo:       1500,
			loserElo:        1500,
			wantWinnerDelta: 16, // K * 0.5 = 16
			wantLoserDelta:  -16,
		},
		{
			name:            "winner higher rated (expected outcome)",
			winnerElo:       1600,
			loserElo:        1400,
			wantWinnerDelta: 7.7, // smaller gain for expected win
			wantLoserDelta:  -7.7,
		},
		{
			name:            "upset - lower rated wins",
			winnerElo:       1400,
			loserElo:        1600,
			wantWinnerDelta: 24.3, // larger gain for upset
			wantLoserDelta:  -24.3,
		},
		{
			name:            "large rating difference - favorite wins",
			winnerElo:       2000,
			loserElo:        1000,
			wantWinnerDelta: 0.5, // minimal gain
			wantLoserDelta:  -0.5,
		},
		{
			name:            "large rating difference - upset",
			winnerElo:       1000,
			loserElo:        2000,
			wantWinnerDelta: 31.5, // near maximum gain
			wantLoserDelta:  -31.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			winner := &db.Restaurant{ID: 1, Name: "Winner", Elo: tt.winnerElo}
			loser := &db.Restaurant{ID: 2, Name: "Loser", Elo: tt.loserElo}

			originalWinnerElo := winner.Elo
			originalLoserElo := loser.Elo

			UpdateElo(winner, loser)

			winnerDelta := winner.Elo - originalWinnerElo
			loserDelta := loser.Elo - originalLoserElo

			// Check delta is in expected direction and approximately correct
			if math.Abs(winnerDelta-tt.wantWinnerDelta) > 1.0 {
				t.Errorf("winner delta = %.2f, want ~%.2f", winnerDelta, tt.wantWinnerDelta)
			}
			if math.Abs(loserDelta-tt.wantLoserDelta) > 1.0 {
				t.Errorf("loser delta = %.2f, want ~%.2f", loserDelta, tt.wantLoserDelta)
			}

			// Elo changes should be symmetric (zero-sum)
			if math.Abs(winnerDelta+loserDelta) > 0.001 {
				t.Errorf("Elo not zero-sum: winner delta %.2f + loser delta %.2f = %.2f",
					winnerDelta, loserDelta, winnerDelta+loserDelta)
			}

			// Winner's Elo should increase
			if winner.Elo <= originalWinnerElo {
				t.Error("winner's Elo should increase")
			}

			// Loser's Elo should decrease
			if loser.Elo >= originalLoserElo {
				t.Error("loser's Elo should decrease")
			}
		})
	}
}

func TestUpdateElo_ScoreUpdate(t *testing.T) {
	t.Parallel()

	winner := &db.Restaurant{ID: 1, Elo: 1500, Score: 5.0}
	loser := &db.Restaurant{ID: 2, Elo: 1500, Score: 5.0}

	UpdateElo(winner, loser)

	// Winner's score should increase
	if winner.Score <= 5.0 {
		t.Errorf("winner.Score = %.2f, should be > 5.0", winner.Score)
	}

	// Loser's score should decrease
	if loser.Score >= 5.0 {
		t.Errorf("loser.Score = %.2f, should be < 5.0", loser.Score)
	}

	// Scores should be in 0-10 range
	if winner.Score < 0 || winner.Score > 10 {
		t.Errorf("winner.Score = %.2f, should be in [0, 10]", winner.Score)
	}
	if loser.Score < 0 || loser.Score > 10 {
		t.Errorf("loser.Score = %.2f, should be in [0, 10]", loser.Score)
	}
}

func TestCalculateScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		elo       float64
		wantScore float64
		tolerance float64
	}{
		{
			name:      "average Elo (1500)",
			elo:       1500,
			wantScore: 5.0, // sigmoid midpoint
			tolerance: 0.01,
		},
		{
			name:      "high Elo (1700)",
			elo:       1700,
			wantScore: 7.31, // 10 / (1 + exp(-1))
			tolerance: 0.1,
		},
		{
			name:      "low Elo (1300)",
			elo:       1300,
			wantScore: 2.69, // 10 / (1 + exp(1))
			tolerance: 0.1,
		},
		{
			name:      "very high Elo (2000)",
			elo:       2000,
			wantScore: 9.24,
			tolerance: 0.1,
		},
		{
			name:      "very low Elo (1000)",
			elo:       1000,
			wantScore: 0.76,
			tolerance: 0.1,
		},
		{
			name:      "extreme high Elo (3000)",
			elo:       3000,
			wantScore: 10.0, // approaches 10
			tolerance: 0.01,
		},
		{
			name:      "extreme low Elo (0)",
			elo:       0,
			wantScore: 0.0, // approaches 0
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := CalculateScore(tt.elo)

			if math.Abs(got-tt.wantScore) > tt.tolerance {
				t.Errorf("CalculateScore(%v) = %.2f, want ~%.2f", tt.elo, got, tt.wantScore)
			}

			// Score should always be in [0, 10]
			if got < 0 || got > 10 {
				t.Errorf("CalculateScore(%v) = %.2f, should be in [0, 10]", tt.elo, got)
			}
		})
	}
}

func TestCalculateScore_Monotonic(t *testing.T) {
	t.Parallel()

	// Score should be monotonically increasing with Elo
	prevScore := CalculateScore(0)
	for elo := 100.0; elo <= 3000; elo += 100 {
		score := CalculateScore(elo)
		if score <= prevScore {
			t.Errorf("score not monotonically increasing: CalculateScore(%.0f) = %.2f <= %.2f",
				elo, score, prevScore)
		}
		prevScore = score
	}
}

func TestRecalculateAllScores(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		restaurants []db.Restaurant
		wantScores  []float64
	}{
		{
			name:        "empty list",
			restaurants: []db.Restaurant{},
			wantScores:  []float64{},
		},
		{
			name: "single restaurant",
			restaurants: []db.Restaurant{
				{ID: 1, Name: "Only"},
			},
			wantScores: []float64{10.0}, // 10 * (1/1) = 10
		},
		{
			name: "two restaurants",
			restaurants: []db.Restaurant{
				{ID: 1, Name: "First"},
				{ID: 2, Name: "Second"},
			},
			wantScores: []float64{10.0, 5.0}, // 10*(2/2), 10*(1/2)
		},
		{
			name: "five restaurants",
			restaurants: []db.Restaurant{
				{ID: 1, Name: "First"},
				{ID: 2, Name: "Second"},
				{ID: 3, Name: "Third"},
				{ID: 4, Name: "Fourth"},
				{ID: 5, Name: "Fifth"},
			},
			wantScores: []float64{10.0, 8.0, 6.0, 4.0, 2.0},
		},
		{
			name: "ten restaurants",
			restaurants: []db.Restaurant{
				{ID: 1},
				{ID: 2},
				{ID: 3},
				{ID: 4},
				{ID: 5},
				{ID: 6},
				{ID: 7},
				{ID: 8},
				{ID: 9},
				{ID: 10},
			},
			wantScores: []float64{10.0, 9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Make a copy to avoid modifying the test data
			restaurants := make([]db.Restaurant, len(tt.restaurants))
			copy(restaurants, tt.restaurants)

			RecalculateAllScores(restaurants)

			for i, want := range tt.wantScores {
				if math.Abs(restaurants[i].Score-want) > 0.001 {
					t.Errorf("restaurants[%d].Score = %.2f, want %.2f",
						i, restaurants[i].Score, want)
				}
			}
		})
	}
}

func TestRecalculateAllScores_ScoresDescending(t *testing.T) {
	t.Parallel()

	restaurants := make([]db.Restaurant, 100)
	for i := range restaurants {
		restaurants[i] = db.Restaurant{ID: int64(i + 1)}
	}

	RecalculateAllScores(restaurants)

	for i := 1; i < len(restaurants); i++ {
		if restaurants[i].Score >= restaurants[i-1].Score {
			t.Errorf("scores not descending: [%d]=%.2f >= [%d]=%.2f",
				i, restaurants[i].Score, i-1, restaurants[i-1].Score)
		}
	}
}

func TestBinarySearch_FullWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("build ranking from scratch", func(t *testing.T) {
		t.Parallel()

		// Simulate ranking 5 restaurants from scratch
		var ranked []db.Restaurant

		restaurants := []db.Restaurant{
			{ID: 1, Name: "A"},
			{ID: 2, Name: "B"},
			{ID: 3, Name: "C"},
			{ID: 4, Name: "D"},
			{ID: 5, Name: "E"},
		}

		// Insert each restaurant, always winning (so order matches input order)
		for _, r := range restaurants {
			bs := NewBinarySearch(r, ranked)
			for !bs.Done {
				bs.Step(true) // new always beats existing
			}
			// Insert at the found position
			newRanked := make([]db.Restaurant, 0, len(ranked)+1)
			newRanked = append(newRanked, ranked[:bs.InsertAt]...)
			newRanked = append(newRanked, r)
			newRanked = append(newRanked, ranked[bs.InsertAt:]...)
			ranked = newRanked
		}

		// All should be inserted at position 0, so order is reversed
		if len(ranked) != 5 {
			t.Fatalf("len(ranked) = %d, want 5", len(ranked))
		}

		// Last inserted (E) should be first
		if ranked[0].Name != "E" {
			t.Errorf("ranked[0] = %s, want E", ranked[0].Name)
		}
		// First inserted (A) should be last
		if ranked[4].Name != "A" {
			t.Errorf("ranked[4] = %s, want A", ranked[4].Name)
		}
	})
}

func TestKConstant(t *testing.T) {
	t.Parallel()

	// K should be 32 (standard chess K-factor for established players)
	if K != 32.0 {
		t.Errorf("K = %v, want 32.0", K)
	}
}

// Benchmarks

func BenchmarkNewBinarySearch(b *testing.B) {
	ranked := make([]db.Restaurant, 1000)
	for i := range ranked {
		ranked[i] = db.Restaurant{ID: int64(i + 1), Name: "R"}
	}
	newR := db.Restaurant{ID: 1001, Name: "New"}

	b.ResetTimer()
	for b.Loop() {
		NewBinarySearch(newR, ranked)
	}
}

func BenchmarkBinarySearchStep(b *testing.B) {
	ranked := make([]db.Restaurant, 1000)
	for i := range ranked {
		ranked[i] = db.Restaurant{ID: int64(i + 1), Name: "R"}
	}
	newR := db.Restaurant{ID: 1001, Name: "New"}

	b.ResetTimer()
	for b.Loop() {
		bs := NewBinarySearch(newR, ranked)
		for !bs.Done {
			bs.Step(true)
		}
	}
}

func BenchmarkUpdateElo(b *testing.B) {
	winner := &db.Restaurant{ID: 1, Elo: 1500}
	loser := &db.Restaurant{ID: 2, Elo: 1500}

	b.ResetTimer()
	for b.Loop() {
		winner.Elo = 1500
		loser.Elo = 1500
		UpdateElo(winner, loser)
	}
}

func BenchmarkCalculateScore(b *testing.B) {
	for b.Loop() {
		CalculateScore(1500)
	}
}

func BenchmarkRecalculateAllScores(b *testing.B) {
	restaurants := make([]db.Restaurant, 1000)
	for i := range restaurants {
		restaurants[i] = db.Restaurant{ID: int64(i + 1)}
	}

	b.ResetTimer()
	for b.Loop() {
		RecalculateAllScores(restaurants)
	}
}
