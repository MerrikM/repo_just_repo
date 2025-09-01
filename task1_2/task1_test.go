package main

import "testing"

func TestGetScore(t *testing.T) {
	stamps := []ScoreStamp{
		{Offset: 0, Score: Score{0, 0}},
		{Offset: 5, Score: Score{1, 0}},
		{Offset: 10, Score: Score{1, 1}},
	}

	tests := []struct {
		name   string
		offset int
		want   Score
	}{
		{
			name:   "Exact match (offset=5)",
			offset: 5,
			want:   Score{1, 0},
		},
		{
			name:   "Between offsets (offset=7)",
			offset: 7,
			want:   Score{1, 0},
		},
		{
			name:   "Offset less than first (offset=-1)",
			offset: -1,
			want:   Score{0, 0},
		},
		{
			name:   "Offset greater than last (offset=20)",
			offset: 20,
			want:   Score{1, 1},
		},
		{
			name:   "Single element slice",
			offset: 100,
			want:   Score{5, 7},
		},
	}

	// отдельный случай для одного элемента
	single := []ScoreStamp{
		{Offset: 0, Score: Score{5, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Score
			if tt.name == "Single element slice" {
				got = getScore(single, tt.offset)
			} else {
				got = getScore(stamps, tt.offset)
			}

			if got != tt.want {
				t.Errorf("offset=%d: expected %v, got %v", tt.offset, tt.want, got)
			}
		})
	}
}
