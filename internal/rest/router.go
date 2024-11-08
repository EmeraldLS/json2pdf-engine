package rest

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/emeraldls/platnova-task/internal/generator"
	"github.com/emeraldls/platnova-task/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/unidoc/unipdf/v3/creator"
)

func generatePDF(w http.ResponseWriter, r *http.Request) {
	var apiKey = os.Getenv("apiKey")
	r.ParseForm()

	customerName := r.FormValue("customerName")
	customerAddressLine1 := r.FormValue("addressLine1")
	customerAddressLine2 := r.FormValue("addressLine2")
	city := r.FormValue("city")
	county := r.FormValue("county")
	postcode := r.FormValue("postcode")

	savingsOpeningBalance, _ := strconv.ParseFloat(r.FormValue("savingsOpeningBalance"), 64)
	savingsMoneyIn, _ := strconv.ParseFloat(r.FormValue("savingsMoneyIn"), 64)
	savingsMoneyOut, _ := strconv.ParseFloat(r.FormValue("savingsMoneyOut"), 64)
	checkingOpeningBalance, _ := strconv.ParseFloat(r.FormValue("checkingOpeningBalance"), 64)
	checkingMoneyIn, _ := strconv.ParseFloat(r.FormValue("checkingMoneyIn"), 64)
	checkingMoneyOut, _ := strconv.ParseFloat(r.FormValue("checkingMoneyOut"), 64)
	iban := r.FormValue("iban")
	bic := r.FormValue("bic")

	transactionDates := r.Form["transactionDate"]
	transactionDescriptions := r.Form["transactionDescription"]
	transactionMoneyIns := r.Form["transactionMoneyIn"]
	transactionMoneyOuts := r.Form["transactionMoneyOut"]

	var transactions []types.AccountTransactions

	for i := 0; i < len(transactionDates); i++ {
		transactionDate := transactionDates[i]
		transactionDescription := transactionDescriptions[i]
		transactionMoneyIn, _ := strconv.ParseFloat(transactionMoneyIns[i], 64)
		transactionMoneyOut, _ := strconv.ParseFloat(transactionMoneyOuts[i], 64)

		balance := savingsOpeningBalance + savingsMoneyIn - savingsMoneyOut + transactionMoneyIn - transactionMoneyOut

		transactions = append(transactions, types.AccountTransactions{
			Date:        transactionDate,
			Description: transactionDescription,
			MoneyIn:     transactionMoneyIn,
			MoneyOut:    transactionMoneyOut,
			Balance:     balance,
		})
	}

	accStmt := types.AccountStatement{
		CustomerName: customerName,
		CustomerAddress: types.CustomerAddress{
			AddressLine1: customerAddressLine1,
			AddressLine2: customerAddressLine2,
			City:         city,
			County:       county,
			Postcode:     postcode,
		},
		BalanceSummary: []types.BalanceSummary{
			{
				Product:        "Savings Account",
				OpeningBalance: savingsOpeningBalance,
				MoneyIn:        savingsMoneyIn,
				MoneyOut:       savingsMoneyOut,
				ClosingBalance: savingsOpeningBalance + savingsMoneyIn - savingsMoneyOut,
			},
			{
				Product:        "Checking Account",
				OpeningBalance: checkingOpeningBalance,
				MoneyIn:        checkingMoneyIn,
				MoneyOut:       checkingMoneyOut,
				ClosingBalance: checkingOpeningBalance + checkingMoneyIn - checkingMoneyOut,
			},
		},
		AccountTransactions: transactions,
		IBANDetails: []types.IBANDetails{
			{
				IBAN: iban,
				BIC:  bic,
			},
		},
	}

	c := types.NewClient(creator.New(), apiKey)
	fn, err := generator.GenerateAccountStatementPDF(*c, accStmt)
	if err != nil {
		http.Error(w, fmt.Sprintf("error generating document: %v", err), http.StatusInternalServerError)
		return
	}

	pdfFile, err := os.Open(fn + ".pdf")
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening generated document: %v", err), http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fn+".pdf")
	w.Header().Set("Content-Type", "application/pdf")

	_, err = io.Copy(w, pdfFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error sending file: %v", err), http.StatusInternalServerError)
		return
	}

	go func() {
		time.Sleep(5 * time.Second)
		err = os.Remove(fn + ".pdf")
		if err != nil {
			slog.Error("error deleting file", "err", err)
		}
	}()
}

func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..")
}

func SetupRoutes() {

	s := gin.Default()

	s.LoadHTMLGlob(filepath.Join(GetProjectRoot(), "templates/*"))
	s.Static("/static", filepath.Join(GetProjectRoot(), "static"))

	s.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	s.POST("/generate", func(ctx *gin.Context) {
		generatePDF(ctx.Writer, ctx.Request)
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "2222"

	}

	s.Run("0.0.0.0:" + PORT)

}
