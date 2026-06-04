package pdf

import (
	"fmt"
	"strconv"
	"time"

	"github.com/diother/hintermann-stripe-docgen/internal/repo"
)

func pagesNeeded(itemsLength int) int {
	const (
		firstPageCapacity      = 8
		subsequentPageCapacity = 12
	)
	remainingItems := itemsLength - firstPageCapacity
	var totalPages int

	if remainingItems > 0 {
		additionalPages := (remainingItems + subsequentPageCapacity - 1) /
			subsequentPageCapacity

		totalPages = 1 + additionalPages
	} else {
		totalPages = 1
	}
	return totalPages
}

func monthDates(year, month int) (string, string, string) {
	periodStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, -1)
	releaseDate := periodStart.AddDate(0, 1, 0)

	return periodStart.Format("2 Jan 2006"),
		periodEnd.Format("2 Jan 2006"),
		releaseDate.Format("2 Jan 2006")
}

func getMonthlyTotals(payouts []repo.Payout) (string, string, string, error) {
	var gross, fee, net int
	for _, p := range payouts {
		g, err := strconv.Atoi(p.Gross)
		if err != nil {
			return "", "", "", err
		}

		f, err := strconv.Atoi(p.Fee)
		if err != nil {
			return "", "", "", err
		}

		n, err := strconv.Atoi(p.Net)
		if err != nil {
			return "", "", "", err
		}

		gross += g
		fee += f
		net += n
	}
	return formatAmount(gross),
		formatAmount(fee),
		formatAmount(net),
		nil
}

func formatAmount(amount int) string {
	return fmt.Sprintf("%.2f lei", float64(amount)/100)
}

func formatStringAmount(amount string) (string, error) {
	a, err := strconv.Atoi(amount)
	if err != nil {
		return "", err
	}

	return formatAmount(a), nil
}
