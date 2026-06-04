package repo

import "path/filepath"

const (
	dataDir     = "data"
	payoutsCsv  = "payouts.csv"
	invoicesCsv = "invoices.csv"
)

func payoutsPath(dir string) string {
	return filepath.Join(dir, payoutsCsv)
}

func invoicesPath(dir string) string {
	return filepath.Join(dir, invoicesCsv)
}
