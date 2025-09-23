package utils

import (
	"math"
	"os"

	"github.com/gagliardetto/solana-go"
)

const LAMPORT_PER_SOL int = 1_000_000_000

func GetPrvKey() (solana.PrivateKey, error) {
	pvk, err := solana.PrivateKeyFromSolanaKeygenFile(os.Getenv("KEYPATH_PATH"))
	if err != nil {
		return solana.PrivateKey{}, err
	}

	return pvk, err
}

func ToLamports(sol_amount float64) int {
	return int(math.Round(sol_amount * float64(LAMPORT_PER_SOL)))
}
