package statstest

import "github.com/bitcoin-sv/spv-wallet/models"

func ExpectedStatsResponse() *models.AdminStats {
	return &models.AdminStats{
		Balance:          0,
		Destinations:     95,
		PaymailAddresses: 20,
		Transactions:     38,
		TransactionsPerDay: map[string]any{
			"20241003": float64(6),
			"20241007": float64(3),
			"20241107": float64(3),
			"20241108": float64(5),
			"20241111": float64(3),
			"20241112": float64(10),
			"20241118": float64(7),
			"20241203": float64(1),
		},
		Utxos:        54,
		UtxosPerType: map[string]any{"pubkeyhash": float64(54)},
		XPubs:        78,
	}
}
