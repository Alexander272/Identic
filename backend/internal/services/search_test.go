package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchService_ValidateTokens(t *testing.T) {
	s := &SearchService{}

	tests := []struct {
		name      string
		reqTokens []string
		pSearch   string
		want      bool
	}{
		{
			name:      "type letter present",
			reqTokens: []string{"прокладка", "снп", "в", "3", "43", "4", "0", "3", "2", "ост26", "260", "454", "99"},
			pSearch:   "прокладка снп в 3 43 4 0 3 2 ост26 260 454 99",
			want:      true,
		},
		{
			name:      "type mismatch rejected",
			reqTokens: []string{"прокладка", "снп", "в", "3", "43", "10", "0", "4", "5", "ост26", "260", "454", "99"},
			pSearch:   "прокладка снп д 3 43 1 6 4 5 304л ост26 260 454 99 50 00 шт",
			want:      false,
		},
		{
			name:      "merged type token accepted",
			reqTokens: []string{"снп", "в", "3", "549", "1", "6", "4", "5"},
			pSearch:   "прокладка снп в3 549 1 6 4 5 575х549х529х4 5",
			want:      true,
		},
		{
			name:      "no type letter in request",
			reqTokens: []string{"прокладка", "снп", "3", "43", "10", "0", "4", "5", "ост26"},
			pSearch:   "прокладка снп д 3 43 10 0 4 5 ост26 260 454 99",
			want:      true,
		},
		{
			name:      "preposition single letter not required",
			reqTokens: []string{"прокладка", "снп", "д", "3", "43", "с", "перемычкой"},
			pSearch:   "прокладка снп д 3 43 перемычкой",
			want:      true,
		},
		{
			name:      "merged token for other letter still rejected",
			reqTokens: []string{"снп", "д", "3", "549"},
			pSearch:   "прокладка снп в3 549 1 6 4 5",
			want:      false,
		},
		{
			name:      "specific token missing rejected",
			reqTokens: []string{"снп", "в", "3", "43"},
			pSearch:   "прокладка снп в 3",
			want:      false,
		},
		{
			name: "long query finds contained short record",
			reqTokens: []string{
				"прокладка", "снп", "в", "3", "43", "1", "6", "4", "5",
				"ост", "26", "260", "454", "99",
				"терморасширенный", "графит", "43", "2", "131кг", "см2", "314", "399", "с",
			},
			pSearch: "прокладка снп в 3 43 1 6 4 5 ост 26 260 454 99",
			want:    true,
		},
		{
			name: "long query rejects garbage candidate",
			reqTokens: []string{
				"прокладка", "снп", "в", "3", "43", "1", "6", "4", "5",
				"ост", "26", "260", "454", "99",
				"терморасширенный", "графит", "43", "2", "131кг", "см2", "314", "399", "с",
			},
			pSearch: "прокладка снп в 2 3 500 25 терморасширенный графит",
			want:    false,
		},
		{
			name: "typo in one specific token tolerated",
			reqTokens: []string{
				"прокладка", "снп", "в", "3", "139", "121", "109", "4", "5",
				"ост26", "260", "454", "99",
				"терморасширенный", "графит", "43", "2", "131кг", "см2", "314", "399", "с",
			},
			pSearch: "прокладка снп в 3 139 141 109 4 5 ост26 260 454 99 терморасширенный графит 43 2 131кг см2 314 399 с",
			want:    true,
		},
		{
			name: "medium query near-size accepted",
			reqTokens: []string{
				"прокладка", "снп", "в", "3", "43", "1", "6", "4", "5", "ост26", "260", "454", "99",
			},
			pSearch: "прокладка снп в 3 43 10 0 4 5 ост26 260 454 99",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, s.validateTokens(tt.pSearch, tt.reqTokens))
		})
	}
}
