package pdf

import (
	"testing"

	"github.com/diother/hintermann-stripe-docgen/internal/repo"
)

func TestPagesNeeded(t *testing.T) {
	testCases := map[string]struct {
		items int
		want  int
	}{
		"empty":               {0, 1},
		"oneItem":             {1, 1},
		"firstPageFull":       {8, 1},
		"oneItemOnSecondPage": {9, 2},
		"secondPageFull":      {20, 2},
		"oneItemOnThirdPage":  {21, 3},
		"thirdPageFull":       {32, 3},
		"oneItemOnFourthPage": {33, 4},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := pagesNeeded(tc.items)

			if got != tc.want {
				t.Errorf(
					"pagesNeeded(%d) = %d, want %d",
					tc.items,
					got,
					tc.want,
				)
			}
		})
	}
}

func TestMonthDates(t *testing.T) {
	inputYear := 2025
	inputMonth := 5

	wantPeriodStart := "1 May 2025"
	wantPeriodEnd := "31 May 2025"
	wantReleaseDate := "1 Jun 2025"

	gotPeriodStart, gotPeriodEnd, gotReleaseDate := monthDates(
		inputYear,
		inputMonth,
	)

	if gotPeriodStart != wantPeriodStart ||
		gotPeriodEnd != wantPeriodEnd ||
		gotReleaseDate != wantReleaseDate {
		t.Errorf(
			"got %v %v %v, want %v %v %v",
			gotPeriodStart, gotPeriodEnd, gotReleaseDate,
			wantPeriodStart, wantPeriodEnd, wantReleaseDate,
		)
	}
}

func TestGetMonthlyTotals(t *testing.T) {
	t.Run("validPayouts", func(t *testing.T) {
		inputPayouts := []repo.Payout{
			{Gross: "1000", Fee: "50", Net: "950"},
			{Gross: "1000", Fee: "50", Net: "950"},
			{Gross: "1000", Fee: "50", Net: "950"},
		}
		wantGross := "30.00 lei"
		wantFee := "1.50 lei"
		wantNet := "28.50 lei"

		gotGross, gotFee, gotNet, err := getMonthlyTotals(inputPayouts)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if gotGross != wantGross ||
			gotFee != wantFee ||
			gotNet != wantNet {
			t.Errorf(
				"got %v %v %v, want %v %v %v",
				gotGross, gotFee, gotNet,
				wantGross, wantFee, wantNet,
			)
		}
	})
	t.Run("invalidPayouts", func(t *testing.T) {
		inputPayouts := []repo.Payout{
			{Gross: "1000", Fee: "not-a-number", Net: "950"},
		}
		_, _, _, err := getMonthlyTotals(inputPayouts)

		if err == nil {
			t.Errorf("Expected error, got none")
		}
	})
}

func TestFormatAmount(t *testing.T) {
	want := "10.00 lei"
	got := formatAmount(1000)

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFormatStringAmount(t *testing.T) {
	t.Run("validString", func(t *testing.T) {
		want := "10.00 lei"
		got, err := formatStringAmount("1000")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("invalidString", func(t *testing.T) {
		_, err := formatStringAmount("not-a-number")

		if err == nil {
			t.Errorf("Expected error, got none")
		}
	})
}
