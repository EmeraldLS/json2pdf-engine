package functions

import (
	"fmt"

	"github.com/emeraldls/platnova-task/internal/types"
	"github.com/emeraldls/platnova-task/internal/utils"
	"github.com/unidoc/unipdf/v3/creator"
)

var cellStyles = map[string]types.CellStyle{
	"heading-left": {
		HAlignment: creator.CellHorizontalAlignmentLeft,
	},
	"heading-right": {

		HAlignment: creator.CellHorizontalAlignmentRight,
	},
	"table-key": {
		Indent:     250,
		HAlignment: creator.CellHorizontalAlignmentLeft,
	},
	"table-value": {
		Indent: 50,
	},
	"note-right": {
		HAlignment: creator.CellHorizontalAlignmentLeft,
		Indent:     50,
	},
	"balance-header": {
		VAlignment: creator.CellVerticalAlignmentMiddle,
	},
	"balance-item": {},
}

func DrawHeading(c *creator.Creator, stmt types.AccountStatement) error {
	table := c.NewTable(2)

	table.SetMargins(0, 0, 0, 0)

	headerStyle := c.NewTextStyle()
	headerStyle.FontSize = 25

	err := utils.DrawCell(table, utils.NewPara(c, "Revolut", headerStyle), cellStyles["heading-left"])
	if err != nil {
		return err
	}

	textStyle := c.NewTextStyle()
	textStyle.FontSize = 10

	var innerRightTable = c.NewTable(1)

	title := utils.NewPara(c, stmt.Title, headerStyle)
	generatedDate := utils.NewPara(c, stmt.GeneratedDate, textStyle)
	bankName := utils.NewPara(c, stmt.BankName, textStyle)

	err = utils.DrawCell(innerRightTable, title, cellStyles["heading-right"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(innerRightTable, generatedDate, cellStyles["heading-right"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(innerRightTable, bankName, cellStyles["heading-right"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, innerRightTable, cellStyles["heading-right"])
	if err != nil {
		return err
	}
	err = c.Draw(table)
	if err != nil {
		return err
	}

	return nil
}

func DrawNameSection(c *creator.Creator, stmt types.AccountStatement) error {

	nameStyle := c.NewTextStyle()
	nameStyle.FontSize = 13

	table := c.NewTable(1)
	table.SetMargins(0, 0, 10, 10)
	err := utils.DrawCell(table, utils.NewPara(c, stmt.CustomerName, nameStyle), types.CellStyle{})
	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	return nil
}

func DrawAddressSection(c *creator.Creator, stmt types.AccountStatement) error {
	table := c.NewTable(1)
	addressStyle := c.NewTextStyle()
	addressStyle.FontSize = 8
	err := utils.DrawCell(table, utils.NewPara(c, stmt.CustomerAddress.AddressLine1, addressStyle), types.CellStyle{})

	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, stmt.CustomerAddress.AddressLine2, addressStyle), types.CellStyle{})

	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, stmt.CustomerAddress.City, addressStyle), types.CellStyle{})

	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, stmt.CustomerAddress.County, addressStyle), types.CellStyle{})

	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, stmt.CustomerAddress.Postcode, addressStyle), types.CellStyle{})

	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	return nil
}

func DrawIBANSection(c *creator.Creator, stmt types.AccountStatement) error {
	for i, detail := range stmt.IBANDetails {

		ibanTable := c.NewTable(2)

		if i == 0 {
			ibanTable.SetMargins(0, 0, -10, 0)
		} else {

			ibanTable.SetMargins(0, 0, 10, 0)
		}

		bicTable := c.NewTable(2)

		keyStyle := c.NewTextStyle()
		keyStyle.FontSize = 8.5
		valueStyle := c.NewTextStyle()
		valueStyle.FontSize = 8.5

		err := utils.DrawCell(ibanTable, utils.NewPara(c, "IBAN", keyStyle), cellStyles["table-key"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(ibanTable, utils.NewPara(c, detail.IBAN, valueStyle), cellStyles["table-value"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(bicTable, utils.NewPara(c, "BIC", keyStyle), cellStyles["table-key"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(bicTable, utils.NewPara(c, detail.BIC, valueStyle), cellStyles["table-value"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(bicTable, utils.NewPara(c, "", valueStyle), types.CellStyle{})
		if err != nil {
			return err
		}

		noteStyle := c.NewTextStyle()
		noteStyle.FontSize = 8.5

		if detail.Note != "" {
			p := c.NewStyledParagraph()
			p.SetEnableWrap(true)
			p.SetText("(" + detail.Note + ")").Style = noteStyle
			err = utils.DrawCell(bicTable, p, cellStyles["note-right"])
			if err != nil {
				return err
			}
		}

		err = c.Draw(ibanTable)
		if err != nil {
			return err
		}

		err = c.Draw(bicTable)
		if err != nil {
			return err
		}
	}
	return nil
}

func DrawBalanceSummary(c *creator.Creator, stmt types.AccountStatement) error {
	table := c.NewTable(1)
	table.SetMargins(0, 0, 10, 0)
	headingStyle := c.NewTextStyle()
	headingStyle.FontSize = 14
	p := utils.NewPara(c, "Balance Summary", headingStyle)

	err := utils.DrawCell(table, p, types.CellStyle{})

	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	table = c.NewTable(5)
	table.SetMargins(0, 0, 10, 0)
	err = table.SetColumnWidths(0.35, 0.2, 0.175, 0.175, 0.1)
	if err != nil {
		return err
	}

	tableHeader := c.NewTextStyle()

	err = utils.DrawCell(table, utils.NewPara(c, "Product", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Opening Balance", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Money Out", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Money In", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	clb := c.NewStyledParagraph()
	clb.SetText("Closing Balance")

	err = utils.DrawCell(table, clb, cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	line := utils.NewLine(c, 0, 0, 0, 0, false, true, 1, creator.ColorBlack, false, nil, 0, 1.0, utils.NewMargins(0, 0, 5, 0))

	err = c.Draw(line)
	if err != nil {
		return err
	}

	tableItem := c.NewTextStyle()

	for _, bs := range stmt.BalanceSummary {
		table = c.NewTable(5)
		table.SetMargins(0, 0, 10, 0)
		err = table.SetColumnWidths(0.35, 0.2, 0.175, 0.175, 0.1)
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, bs.Product, tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", bs.OpeningBalance), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", bs.MoneyOut), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", bs.MoneyIn), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", bs.ClosingBalance), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = c.Draw(table)
		if err != nil {
			return err
		}

		line = utils.NewLine(c, 0, 0, 0, 0, false, true, 1, creator.ColorBlack, false, nil, 0, 1.0, utils.NewMargins(0, 0, 5, 0))

		err = c.Draw(line)
		if err != nil {
			return err
		}

	}

	table = c.NewTable(5)
	table.SetMargins(0, 0, 10, 0)
	err = table.SetColumnWidths(0.35, 0.2, 0.175, 0.175, 0.1)
	if err != nil {
		return err
	}
	tableItem = c.NewTextStyle()

	err = utils.DrawCell(table, utils.NewPara(c, "Total", tableItem), cellStyles["balance-item"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", stmt.GetTotalOpeningBalanceSummary()), tableItem), cellStyles["balance-item"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", stmt.GetTotalMoneyOutBalanceSummary()), tableItem), cellStyles["balance-item"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", stmt.GetTotalMoneyInBalanceSummary()), tableItem), cellStyles["balance-item"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", stmt.GetTotalClosingBalanceSummary()), tableItem), cellStyles["balance-item"])
	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	return nil
}

func DrawAccountTransactionsSummary(c *creator.Creator, stmt types.AccountStatement) error {
	table := c.NewTable(1)
	table.SetMargins(0, 0, 50, 0)
	headingStyle := c.NewTextStyle()
	headingStyle.FontSize = 14
	p := utils.NewPara(c, "Account Transactions from 1 February 2023 to 29 March 2023", headingStyle)

	err := utils.DrawCell(table, p, types.CellStyle{})

	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	table = c.NewTable(5)
	table.SetMargins(0, 0, 10, 0)
	err = table.SetColumnWidths(0.2, 0.35, 0.175, 0.175, 0.1)
	if err != nil {
		return err
	}

	tableHeader := c.NewTextStyle()

	err = utils.DrawCell(table, utils.NewPara(c, "Date", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Description", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Money Out", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Money In", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = utils.DrawCell(table, utils.NewPara(c, "Balance", tableHeader), cellStyles["balance-header"])
	if err != nil {
		return err
	}

	err = c.Draw(table)
	if err != nil {
		return err
	}

	line := utils.NewLine(c, 0, 0, 0, 0, false, true, 1, creator.ColorBlack, false, nil, 0, 1.0, utils.NewMargins(0, 0, 5, 0))

	err = c.Draw(line)
	if err != nil {
		return err
	}

	tableItem := c.NewTextStyle()

	for _, tx := range stmt.AccountTransactions {
		table = c.NewTable(5)
		table.SetMargins(0, 0, 10, 0)
		err = table.SetColumnWidths(0.2, 0.35, 0.175, 0.175, 0.1)
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, tx.Date, tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, tx.Description, tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", tx.MoneyOut), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", tx.MoneyIn), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = utils.DrawCell(table, utils.NewPara(c, fmt.Sprintf("$%.2f", tx.Balance), tableItem), cellStyles["balance-item"])
		if err != nil {
			return err
		}

		err = c.Draw(table)
		if err != nil {
			return err
		}

		line := utils.NewLine(c, 0, 0, 0, 0, false, true, 1, creator.ColorBlack, false, nil, 0, 1.0, utils.NewMargins(0, 0, 5, 0))

		err = c.Draw(line)
		if err != nil {
			return err
		}

	}

	return nil
}

func DrawFooter(c *creator.Creator) {
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		img, err := c.NewImageFromFile("./qrcode.png")
		if err != nil {
			panic(err)
		}
		img.ScaleToWidth(40)

		table := c.NewTable(3)
		table.SetColumnWidths(0.1, 0.2, 0.7)
		table.SetMargins(55, 55, 0, 0)

		_ = utils.DrawCell(table, img, types.CellStyle{})

		headingStyle := c.NewTextStyle()
		headingStyle.FontSize = 6

		captionStyle := c.NewTextStyle()
		captionStyle.FontSize = 4

		stylePara := c.NewStyledParagraph()
		stylePara.EnableWordWrap(true)
		stylePara.Append("Report Lost or Stolen Card").Style = headingStyle
		stylePara.Append("+370 5 214 3608").Style = captionStyle
		stylePara.Append("Get Help Directly In App").Style = headingStyle
		stylePara.Append("Scan the QR code ").Style = captionStyle

		_ = utils.DrawCell(table, stylePara, types.CellStyle{})

		p := c.NewParagraph(`Revolut Bank UAB, Konstitucijos Ave. 21B, 08130 Vilnius, Lithuania, company number 304580906. Revolut Bank UAB is an electronic money institution authorized by the Bank of Lithuania, authorization code LB001994. Revolut Bank UAB is a member of the State Company “Deposit and Investment Insurance” (Valstybės įmonė “Indėlių ir investicijų draudimas”), and deposits are insured according to the conditions established by the laws of the Republic of Lithuania. Your funds are protected by the deposit guarantee scheme up to €100,000 per depositor. Additional protection for deposits exceeding €100,000 is subject to the applicable terms and conditions.
		`)
		p.SetFontSize(5)

		_ = utils.DrawCell(table, p, types.CellStyle{})

		_ = block.Draw(table)

		p = c.NewParagraph("2023 Revolut Bank UAB")
		p.SetFontSize(8)
		p.SetPos(60, 40)
		block.Draw(p)

		strPage := fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages)
		p = c.NewParagraph(strPage)

		p.SetFontSize(8)
		p.SetPos(520, 40)
		block.Draw(p)
	})
}
