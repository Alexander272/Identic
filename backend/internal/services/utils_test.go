package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeString_FoldLatinToCyrillic(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty",
			in:   "",
			want: "",
		},
		{
			name: "latin look-alikes folded to cyrillic",
			in:   "ПРОКЛАДКА СНП-В-3-139-121-109-4,5 ОСТ26.260.454-99+ТЕРМОРАСШИРЕННЫЙ ГРАФИТ 43,2-131КГ/СМ2 314-399°С",
			want: "прокладка снп в 3 139 121 109 4 5 ост26 260 454 99 терморасширенный графит 43 2 131кг см2 314 399 с",
		},
		{
			name: "all-latin input folded to cyrillic",
			in:   "ПPOKЛAДKA CHП-B-3-139-121-109-4,5 OCT26.260.454-99+TEPMOPACШИPEHHЫЙ ГPAФИT 43,2-131KГ/CM2 314-399°C",
			want: "прокладка снп в 3 139 121 109 4 5 ост26 260 454 99 терморасширенный графит 43 2 131кг см2 314 399 с",
		},
		{
			name: "latin letters in mixed text",
			in:   "Прокладка СНП B 3 44x31x25 4,5 110",
			want: "прокладка снп в 3 44х31х25 4 5 110",
		},
		{
			name: "punct removed and spaces collapsed",
			in:   "ГРАФИТ 43,2-131КГ/СМ2  314-399°С",
			want: "графит 43 2 131кг см2 314 399 с",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeString(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
