package main

import (
	"math/big"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/utils/starknet"
	"fmt"
)	

func main() {
	// First hash calculation
	hash := starknet.GetLimitOrderMsg(
		big.NewInt(0x1),
		big.NewInt(0x2),
		true,
		big.NewInt(0x3),
		big.NewInt(0x4),
		big.NewInt(0x5),
		big.NewInt(0x6),
		123456789,
		100000,
		1000000,
	)
	
	fmt.Println("First hash result:", hash.String())

	// Second hash calculation
	hash = starknet.GetLimitOrderMsg(
		big.NewInt(0x2), // Changed from 0x1
		big.NewInt(0x3), // Changed from 0x2
		false, // Changed from true
		big.NewInt(0x4), // Changed from 0x3
		big.NewInt(0x5), // Changed from 0x4
		big.NewInt(0x6), // Changed from 0x5
		big.NewInt(0x7), // Changed from 0x6
		987654321, // Changed from 123456789
		200000, // Changed from 100000
		2000000, // Changed from 1000000
	)
	
	fmt.Println("Second hash result:", hash.String())

	// Third hash calculation
	hash = starknet.GetLimitOrderMsg(
		big.NewInt(0x3), // Changed from 0x1
		big.NewInt(0x4), // Changed from 0x2
		true,
		big.NewInt(0x5), // Changed from 0x3
		big.NewInt(0x6), // Changed from 0x4
		big.NewInt(0x7), // Changed from 0x5
		big.NewInt(0x8), // Changed from 0x6
		456789123, // Changed from 123456789
		300000, // Changed from 100000
		3000000, // Changed from 1000000
	)
	fmt.Println("Third hash result:", hash.String())
}
