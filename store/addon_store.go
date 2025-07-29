package store

import "sync"

type AddOnStats struct {
	TotalSales         int
	TotalRevenueCents  int
	TotalCommissionCents int
}

var stats AddOnStats
var mutex sync.Mutex

// InitStore initializes the store package (add any setup logic here)
func InitStore() {
    // Initialization logic, if needed
}

func RecordAddOnPurchase(amountCents int) {
	mutex.Lock()
	defer mutex.Unlock()

	stats.TotalSales++
	stats.TotalRevenueCents += amountCents
	stats.TotalCommissionCents += amountCents * 20 / 100
}

func GetStats() AddOnStats {
	mutex.Lock()
	defer mutex.Unlock()

	return stats
}
