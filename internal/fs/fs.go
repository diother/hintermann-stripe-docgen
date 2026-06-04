package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/signintech/gopdf"
)

func PrepareBundle(year, month int) (string, string, error) {
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return "", "", err
	}

	monthlyRoot := filepath.Join(
		distDir,
		fmt.Sprintf("%d_%02d", year, month),
	)

	if err := os.RemoveAll(monthlyRoot); err != nil {
		return "", "", err
	}

	if err := os.MkdirAll(monthlyRoot, 0755); err != nil {
		return "", "", err
	}

	payoutRoot := filepath.Join(monthlyRoot, payoutsDir)
	if err := os.MkdirAll(payoutRoot, 0755); err != nil {
		return "", "", err
	}

	return monthlyRoot, payoutRoot, nil
}

func WriteMonthly(root string, pdf *gopdf.GoPdf) error {
	path := filepath.Join(root, monthlyReportPdf)
	return pdf.WritePdf(path)
}

func WritePayout(
	payoutRoot string,
	payoutId string,
	pdf *gopdf.GoPdf,
) (
	string,
	error,
) {
	currDir := filepath.Join(payoutRoot, payoutId)
	if err := os.MkdirAll(currDir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(currDir, payoutReportPdf)
	if err := pdf.WritePdf(path); err != nil {
		return "", nil
	}

	invoiceRoot := filepath.Join(currDir, invoicesDir)
	if err := os.MkdirAll(invoiceRoot, 0755); err != nil {
		return "", err
	}

	return invoiceRoot, nil
}

func WriteInvoice(
	invoiceRoot string,
	invoiceId string,
	pdf *gopdf.GoPdf,
) error {
	path := filepath.Join(invoiceRoot, invoiceId+".pdf")
	return pdf.WritePdf(path)
}
