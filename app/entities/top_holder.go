package entities

import "time"

type TopHolder struct {
	Mint                   string
	Top5Wallets            float64
	Top10Wallets           float64
	Top20Wallets           float64
	HasSingleLargestHolder bool
	TrackedAt              time.Time
}

func NewTopHolder(
	mint string,
	top_5_wallets,
	top_10_wallets,
	top_20_wallets float64,
	has_single_largest_holder bool,
) TopHolder {
	return TopHolder{
		Mint:                   mint,
		Top5Wallets:            top_5_wallets,
		Top10Wallets:           top_10_wallets,
		Top20Wallets:           top_20_wallets,
		HasSingleLargestHolder: has_single_largest_holder,
	}
}
