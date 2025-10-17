package starknet

import (
	"fmt"
	"math/big"

	felt "github.com/NethermindEth/juno/core/felt"
)

// bigIntToFelt converts a big.Int to a felt.Felt, returning an error if the value
// is negative or greater than or equal to the Stark field prime. No modulo is applied.
func bigIntToFelt(x *big.Int) (*felt.Felt, error) {
	if x == nil {
		return nil, fmt.Errorf("nil big.Int")
	}
	p, _ := new(big.Int).SetString("0x800000000000011000000000000000000000000000000000000000000000001", 0)
	if x.Sign() < 0 || x.Cmp(p) >= 0 {
		return nil, fmt.Errorf("value %s out of field range", x.String())
	}
	f := new(felt.Felt)
	f.SetBigInt(x)
	return f, nil
}
