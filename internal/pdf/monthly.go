package pdf

import (
	"fmt"

	"github.com/diother/hintermann-stripe-docgen/internal/repo"
	"github.com/signintech/gopdf"
)

func GenerateMonthly(
	year,
	month int,
	payouts []repo.Payout,
) (
	*gopdf.GoPdf,
	error,
) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err := setFonts(pdf); err != nil {
		return nil, fmt.Errorf("failed setting fonts: %w", err)
	}
	resetTextStyles(pdf)

	itemsLength := len(payouts)
	pagesNeeded := pagesNeeded(itemsLength)
	currentPage := 1

	periodStart, periodEnd, releaseDate := monthDates(year, month)

	if err := addMonthlyReportHeader(pdf, releaseDate); err != nil {
		return nil, fmt.Errorf("failed adding the header: %w", err)
	}

	err := addMonthlyReportFooter(pdf, currentPage, pagesNeeded)
	if err != nil {
		return nil, fmt.Errorf("failed adding the footer: %w", err)
	}

	gross, fee, net, err := getMonthlyTotals(payouts)
	if err != nil {
		return nil, fmt.Errorf("failed getting totals: %w", err)
	}

	addMonthlyPayoutSummary(pdf, gross, fee, net, periodStart, periodEnd)
	addMonthlyPayoutTable(pdf, firstPageTableY)

	currentY := firstPageStartY
	maxItemsPerPage := firstPageCapacity

	var itemCounter int
	for _, payout := range payouts {
		if itemCounter == maxItemsPerPage {
			pdf.AddPage()
			currentPage++

			err := addMonthlyReportSecondaryHeader(pdf)
			if err != nil {
				return nil, fmt.Errorf(
					"failed adding the secondary header: %w",
					err,
				)
			}

			err = addMonthlyReportFooter(pdf, currentPage, pagesNeeded)
			if err != nil {
				return nil, fmt.Errorf("failed adding the footer: %w", err)
			}

			addMonthlyPayoutTable(pdf, subsequentPageTableY)

			currentY = subsequentPageStartY
			itemCounter = 0
			maxItemsPerPage = subsequentPageCapacity
		}

		err := addMonthlyPayoutProduct(pdf, payout, currentY)
		if err != nil {
			return nil, fmt.Errorf("failed adding the payout product: %w", err)
		}

		currentY += itemHeight
		itemCounter++
	}

	return pdf, nil
}

func addMonthlyReportHeader(pdf *gopdf.GoPdf, releaseDate string) error {
	const startY = marginTop

	err := addImage(pdf, stripeLogoPath, marginLeft, startY, 51, 21)
	if err != nil {
		return err
	}

	setText(pdf, marginLeft, startY+31, "Stripe Payments Europe, Limited")
	setText(pdf, marginLeft, startY+47, "The One Building")
	setText(pdf, marginLeft, startY+63, "1 Grand Canal Street Lower")
	setText(pdf, marginLeft, startY+79, "Dublin 2")
	setText(pdf, marginLeft, startY+95, "Co. Dublin")
	setText(pdf, marginLeft, startY+111, "Ireland")

	setText(pdf, 312, startY+31, "Data emiterii:")
	setRightAlignedText(pdf, marginRight, startY+31, releaseDate)
	setText(pdf, 312, startY+47, "Nr. cont:")
	setRightAlignedText(pdf, marginRight, startY+47, "acct_1PVfUvDXCtuWOFq8")
	setText(pdf, 312, startY+63, "Proprietar cont:")
	setRightAlignedText(
		pdf,
		marginRight,
		startY+63,
		"Asociația de Caritate Hintermann",
	)
	setText(pdf, 312, startY+79, "Adresă:")
	setRightAlignedText(pdf, marginRight, startY+79, "Strada Spicului, Nr. 12")
	setRightAlignedText(pdf, marginRight, startY+95, "Bl. 40, Sc. A, Ap. 12")
	setRightAlignedText(pdf, marginRight, startY+111, "Brașov, România")
	setRightAlignedText(pdf, marginRight, startY+127, "500460")

	pdf.SetFont(robotoBold, "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras lunar")

	resetTextStyles(pdf)
	return nil
}

func addMonthlyReportSecondaryHeader(pdf *gopdf.GoPdf) error {
	const startY = marginTop

	err := addImage(pdf, stripeLogoPath, marginLeft, startY, 51, 21)
	if err != nil {
		return err
	}

	pdf.SetFont(robotoBold, "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras lunar")

	resetTextStyles(pdf)
	return nil
}

func addMonthlyReportFooter(
	pdf *gopdf.GoPdf,
	currentPage,
	pagesNeeded int,
) error {
	const endY = marginBottom

	err := addImage(pdf, stripeLogoSmallPath, marginLeft, endY-17, 41, 17)
	if err != nil {
		return err
	}

	pageInfo := fmt.Sprintf("Pagina %d din %d", currentPage, pagesNeeded)
	setText(pdf, 492, endY-15.5, pageInfo)

	pdf.Line(marginLeft, endY-37, marginRight, endY-37)
	return nil
}

func addMonthlyPayoutSummary(
	pdf *gopdf.GoPdf,
	gross,
	fee,
	net,
	periodStart,
	periodEnd string,
) {
	const startY = 211

	period := periodStart + " - " + periodEnd

	setText(pdf, marginLeft, startY+26, period)

	setText(pdf, 312, startY+10, "Preț brut:")
	setText(pdf, 312, startY+26, "Taxe Stripe:")

	setRightAlignedText(pdf, marginRight, startY+10, gross)
	setRightAlignedText(pdf, marginRight, startY+26, "-"+fee)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY+10, "Periodă extras:")

	pdf.SetFont(robotoBold, "", 10)
	setText(pdf, 312, startY+42, "Total:")
	setRightAlignedText(pdf, marginRight, startY+42, net)

	resetTextStyles(pdf)

	pdf.Line(marginLeft, startY-.5, marginRight, startY-.5)
	pdf.Line(marginLeft, startY+63.5, marginRight, startY+63.5)
	pdf.Line(297.5, startY-.5, 298.5, startY+63.5)
}

func addMonthlyPayoutTable(pdf *gopdf.GoPdf, startY float64) {
	setText(pdf, marginLeft, startY, "Plată")
	setText(pdf, 328, startY, "Preț brut")
	setText(pdf, 424.5, startY, "Taxă Stripe")
	setText(pdf, 532, startY, "Total")

	pdf.Line(marginLeft, startY+21.5, marginRight, startY+21.5)
}

func addMonthlyPayoutProduct(
	pdf *gopdf.GoPdf,
	payout repo.Payout,
	startY float64,
) error {
	setText(pdf, marginLeft, startY+16, payout.Created)

	gross, err := formatStringAmount(payout.Gross)
	if err != nil {
		return err
	}
	fee, err := formatStringAmount(payout.Fee)
	if err != nil {
		return err
	}
	net, err := formatStringAmount(payout.Net)
	if err != nil {
		return err
	}

	setRightAlignedText(pdf, 367, startY, gross)
	setRightAlignedText(pdf, 474, startY, "-"+fee)
	setRightAlignedText(pdf, marginRight, startY, net)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY, payout.Id)
	pdf.SetTextColor(94, 100, 112)

	return nil
}
