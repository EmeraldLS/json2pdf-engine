package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/emeraldls/platnova-task/internal/generator"
	"github.com/emeraldls/platnova-task/internal/types"
	"github.com/emeraldls/platnova-task/internal/utils"
	"github.com/joho/godotenv"
	"github.com/unidoc/unipdf/v3/creator"
)

func TestReadJSONFile(t *testing.T) {
	t.Run("json_reader", func(t *testing.T) {
		stmt, err := utils.ReadJSONFile("account_statement.json")
		if err != nil {
			t.Errorf("an error occured reading file: %v ", err)
			return
		}
		t.Logf("Statement: %v ", stmt)
	})
}

func TestPDFGenerator(t *testing.T) {

	godotenv.Load(".env")

	var uniDocAPIKey = os.Getenv("apiKey")
	t.Run("pdf_generator", func(t *testing.T) {

		stmt, err := utils.ReadJSONFile("account_statement.json")
		if err != nil {
			t.Errorf("unable to read file: %v", err)
			return
		}

		fmt.Println("ApiKey ", uniDocAPIKey)

		c := types.NewClient(creator.New(), uniDocAPIKey)

		fn, err := generator.GenerateAccountStatementPDF(*c, *stmt)
		if err != nil {
			t.Errorf("error occured generating pdf: %v", err)
			return
		}

		t.Logf("PDF generated successfully: %s.pdf", fn)
	})
}
