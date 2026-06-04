package repo

import "time"

func getMonthPeriod(year, month int) (time.Time, time.Time) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	return start, end
}

func filterPayoutsByPeriod(
	payouts []Payout,
	year,
	month int,
) (
	[]Payout,
	error,
) {
	start, end := getMonthPeriod(year, month)

	var filtered []Payout
	for _, p := range payouts {
		created, err := time.Parse("2 Jan 2006", p.Created)
		if err != nil {
			return nil, err
		}

		if (created.Equal(start) || created.After(start)) &&
			(created.Equal(end) || created.Before(end)) {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

func filterInvoicesByPayoutId(
	invoices []Invoice,
	payoutId string,
) []Invoice {
	var filtered []Invoice
	for _, i := range invoices {
		if i.PayoutId == payoutId {
			filtered = append(filtered, i)
		}
	}

	return filtered
}
