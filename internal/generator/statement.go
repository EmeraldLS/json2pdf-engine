package generator

import (
	"github.com/emeraldls/platnova-task/internal/functions"
	"github.com/emeraldls/platnova-task/internal/types"
	"github.com/unidoc/unipdf/v3/common/license"
)

func GenerateAccountStatementPDF(c types.Client, stmt types.AccountStatement) (string, error) {

	err := license.SetMeteredKey(c.UniDocAPIKey)
	if err != nil {
		return "", err
	}

	err = functions.DrawHeading(c.Creator, stmt)
	if err != nil {
		return "", err
	}

	err = functions.DrawNameSection(c.Creator, stmt)
	if err != nil {
		return "", err
	}

	if err = functions.DrawAddressSection(c.Creator, stmt); err != nil {
		return "", err
	}

	err = functions.DrawIBANSection(c.Creator, stmt)
	if err != nil {
		return "", err
	}

	if err = functions.DrawBalanceSummary(c.Creator, stmt); err != nil {
		return "", err
	}

	if err = functions.DrawAccountTransactionsSummary(c.Creator, stmt); err != nil {
		return "", err
	}

	functions.DrawFooter(c.Creator)

	fn, err := c.Save()
	if err != nil {
		return "", err
	}

	return fn, nil
}
