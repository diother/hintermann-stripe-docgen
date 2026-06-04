package repo

import (
	"testing"
	"time"
)

func TestGetMonthPeriod(t *testing.T) {
	gotStart, gotEnd := getMonthPeriod(2025, 5)

	wantStart := time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, time.May, 31, 0, 0, 0, 0, time.UTC)

	if !gotStart.Equal(wantStart) {
		t.Errorf("start = %v, want %v", gotStart, wantStart)
	}

	if !gotEnd.Equal(wantEnd) {
		t.Errorf("end = %v, want %v", gotEnd, wantEnd)
	}
}

func TestFilterPayoutsByPeriod(t *testing.T) {
	inputYear := 2025
	inputMonth := 5

	testCases := map[string]struct {
		inputPayouts []Payout
		wantPayouts  []Payout
		expectedErr  bool
	}{
		"insideMonth": {
			[]Payout{{Created: "15 May 2025"}},
			[]Payout{{Created: "15 May 2025"}},
			false,
		},
		"outsideMonth": {
			[]Payout{{Created: "12 Apr 2025"}},
			[]Payout{},
			false,
		},
		"firstDayMonth": {
			[]Payout{{Created: "1 May 2025"}},
			[]Payout{{Created: "1 May 2025"}},
			false,
		},
		"lastDayMonth": {
			[]Payout{{Created: "31 May 2025"}},
			[]Payout{{Created: "31 May 2025"}},
			false,
		},
		"invalidDate": {
			[]Payout{{Created: "not-a-date"}},
			nil,
			true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotPayouts, err := filterPayoutsByPeriod(
				tc.inputPayouts,
				inputYear,
				inputMonth,
			)

			if !tc.expectedErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if tc.expectedErr && err == nil {
				t.Errorf("Expected error, got none")
			}

			for i := range gotPayouts {
				if !equalPayout(gotPayouts[i], tc.wantPayouts[i]) {
					t.Errorf(
						"mismatch at index %d:\ngot  %+v\nwant %+v",
						i,
						gotPayouts[i],
						tc.wantPayouts[i],
					)
				}
			}
		})
	}
}

func TestFilterInvoicesByPayoutId(t *testing.T) {
	inputPayoutId := "po_1"
	inputInvoices := []Invoice{
		{Id: "txn_1", PayoutId: "po_3"},
		{Id: "txn_2", PayoutId: "po_1"},
		{Id: "txn_3", PayoutId: "po_3"},
	}

	wantInvoices := []Invoice{
		{Id: "txn_2", PayoutId: "po_1"},
	}

	gotInvoices := filterInvoicesByPayoutId(inputInvoices, inputPayoutId)

	for i := range gotInvoices {
		if !equalInvoice(gotInvoices[i], wantInvoices[i]) {
			t.Errorf(
				"mismatch at index %d:\ngot  %+v\nwant %+v",
				i,
				gotInvoices[i],
				wantInvoices[i],
			)
		}
	}
}

func equalPayout(a, b Payout) bool {
	return a.Id == b.Id &&
		a.Created == b.Created &&
		a.Gross == b.Gross &&
		a.Fee == b.Fee &&
		a.Net == b.Net
}

func equalInvoice(a, b Invoice) bool {
	return a.Id == b.Id &&
		a.Created == b.Created &&
		a.ClientName == b.ClientName &&
		a.ClientEmail == b.ClientEmail &&
		a.PayoutId == b.PayoutId &&
		a.Gross == b.Gross &&
		a.Fee == b.Fee &&
		a.Net == b.Net
}
