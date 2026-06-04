package repo

import (
	"encoding/csv"
	"fmt"
	"os"
)

func GetPayoutsForMonth(year, month int) ([]Payout, error) {
	payouts, err := readPayouts(dataDir)
	if err != nil {
		return nil, err
	}

	filtered, err := filterPayoutsByPeriod(payouts, year, month)
	if err != nil {
		return nil, err
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no payouts found for %d-%02d", year, month)
	}

	return filtered, nil
}

func GetInvoicesForPayout(payoutId string) ([]Invoice, error) {
	invoices, err := readInvoices(dataDir)
	if err != nil {
		return nil, err
	}

	filtered := filterInvoicesByPayoutId(invoices, payoutId)

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no invoices found for %s", payoutId)
	}

	return filtered, nil
}

func readPayouts(dir string) ([]Payout, error) {
	path := payoutsPath(dir)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Payout{}, nil
		}

		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	payouts := make([]Payout, 0, len(records))

	for _, row := range records {
		if len(row) != 5 {
			return nil, fmt.Errorf("invalid payout row length")
		}
		payouts = append(payouts, Payout{
			Id:      row[0],
			Created: row[1],
			Gross:   row[2],
			Fee:     row[3],
			Net:     row[4],
		})
	}

	return payouts, nil
}

func readInvoices(dir string) ([]Invoice, error) {
	path := invoicesPath(dir)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Invoice{}, nil
		}

		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()

	invoices := make([]Invoice, 0, len(records))

	for _, row := range records {
		if len(row) != 8 {
			return nil, fmt.Errorf("invalid invoice row length")
		}

		invoices = append(invoices, Invoice{
			Id:          row[0],
			Created:     row[1],
			ClientName:  row[2],
			ClientEmail: row[3],
			PayoutId:    row[4],
			Gross:       row[5],
			Fee:         row[6],
			Net:         row[7],
		})
	}

	return invoices, nil
}
