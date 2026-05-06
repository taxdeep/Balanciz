package pages

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"balanciz/internal/models"
)

func TestPayBillsPageAutoLooksUpForeignExchangeRate(t *testing.T) {
	due := time.Date(2026, 6, 4, 0, 0, 0, 0, time.UTC)
	vm := PayBillsVM{
		HasCompany:   true,
		BaseCurrency: "CAD",
		EntryDate:    "2026-06-05",
		OpenBills: []models.Bill{
			{
				ID:           17,
				BillNumber:   "BILL017",
				Vendor:       models.Vendor{Name: "USD Vendor"},
				DueDate:      &due,
				Amount:       decimal.NewFromInt(100),
				BalanceDue:   decimal.NewFromInt(100),
				CurrencyCode: "USD",
			},
		},
		BillAmounts: map[string]string{"17": "100.00"},
	}

	var sb strings.Builder
	if err := PayBills(vm).Render(context.Background(), &sb); err != nil {
		t.Fatalf("render PayBills: %v", err)
	}
	html := sb.String()

	for _, want := range []string{
		`@input.debounce.300ms="onPaymentDateChange()"`,
		`@change="onPaymentDateChange()"`,
		`x-model="exchangeRate"`,
		`@input="onExchangeRateInput()"`,
		`placeholder="auto-lookup"`,
		`The system will auto-fill a rate when available; you can override it here.`,
		`x-text="exchangeRateHint"`,
		`lookupExchangeRate: async function`,
		`/api/exchange-rate?`,
		`transaction_currency_code`,
		`allow_provider_fetch`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected Pay Bills HTML to contain %q", want)
		}
	}
}
