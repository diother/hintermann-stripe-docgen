package app

import (
	"fmt"

	"github.com/diother/hintermann-stripe-docgen/internal/fs"
	"github.com/diother/hintermann-stripe-docgen/internal/pdf"
	"github.com/diother/hintermann-stripe-docgen/internal/repo"
)

func Run(year, month int) error {
	payouts, err := repo.GetPayoutsForMonth(year, month)
	if err != nil {
		return fmt.Errorf("failed to retrieve payouts %w", err)
	}

	monthlyRoot, payoutRoot, err := fs.PrepareBundle(year, month)
	if err != nil {
		return fmt.Errorf("failed to prepare bundle: %w", err)
	}

	monthlyPdf, err := pdf.GenerateMonthly(year, month, payouts)
	if err != nil {
		return fmt.Errorf("failed to render monthly: %w", err)
	}

	if err := fs.WriteMonthly(monthlyRoot, monthlyPdf); err != nil {
		return fmt.Errorf("failed to write monthly: %w", err)
	}

	for _, payout := range payouts {
		invoices, err := repo.GetInvoicesForPayout(payout.Id)
		if err != nil {
			return fmt.Errorf("failed to retrieve invoices %w", err)
		}

		payoutPdf, err := pdf.GeneratePayout(payout, invoices)
		if err != nil {
			return fmt.Errorf("failed to render payout: %w", err)
		}

		invoiceRoot, err := fs.WritePayout(payoutRoot, payout.Id, payoutPdf)
		if err != nil {
			return fmt.Errorf("failed to write payout: %w", err)
		}

		for _, invoice := range invoices {
			invoicePdf, err := pdf.GenerateInvoice(invoice)
			if err != nil {
				return fmt.Errorf("failed to render invoice: %w", err)
			}

			err = fs.WriteInvoice(invoiceRoot, invoice.Id, invoicePdf)
			if err != nil {
				return fmt.Errorf("failed to write invoice: %w", err)
			}
		}
	}

	return nil
}
