package main

import (
	"fmt"
	"math/big"
	"os"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Get private key from environment
	privateKeyHex := os.Getenv("X10_PRIVATE_KEY")
	if privateKeyHex == "" {
		fmt.Println("ERROR: X10_PRIVATE_KEY not set in .env")
		return
	}

	// Parse private key
	privateKey, ok := new(big.Int).SetString(privateKeyHex, 0)
	if !ok {
		fmt.Printf("ERROR: Invalid private key format: %s\n", privateKeyHex)
		return
	}

	fmt.Println("=== Starknet Signing Test ===")
	fmt.Println()

	// Test Case 1: Simple known hash
	testHash1 := "1234567890"
	hash1BigInt, _ := new(big.Int).SetString(testHash1, 10)
	hash1Felt := new(felt.Felt).SetBigInt(hash1BigInt)

	fmt.Printf("Test 1:\n")
	fmt.Printf("  Hash (decimal): %s\n", testHash1)
	fmt.Printf("  Hash (hex): 0x%x\n", hash1BigInt)

	privateKeyFelt := new(felt.Felt).SetBigInt(privateKey)
	r1, s1, err := curve.SignFelts(hash1Felt, privateKeyFelt)
	if err != nil {
		fmt.Printf("  ERROR signing: %v\n", err)
	} else {
		r1BigInt := new(big.Int)
		s1BigInt := new(big.Int)
		r1.BigInt(r1BigInt)
		s1.BigInt(s1BigInt)

		fmt.Printf("  Signature r (decimal): %s\n", r1BigInt.String())
		fmt.Printf("  Signature s (decimal): %s\n", s1BigInt.String())
		fmt.Printf("  Signature r (hex): 0x%x\n", r1BigInt)
		fmt.Printf("  Signature s (hex): 0x%x\n", s1BigInt)
	}
	fmt.Println()

	// Test Case 2: Another known hash
	testHash2 := "9876543210"
	hash2BigInt, _ := new(big.Int).SetString(testHash2, 10)
	hash2Felt := new(felt.Felt).SetBigInt(hash2BigInt)

	fmt.Printf("Test 2:\n")
	fmt.Printf("  Hash (decimal): %s\n", testHash2)
	fmt.Printf("  Hash (hex): 0x%x\n", hash2BigInt)

	r2, s2, err := curve.SignFelts(hash2Felt, privateKeyFelt)
	if err != nil {
		fmt.Printf("  ERROR signing: %v\n", err)
	} else {
		r2BigInt := new(big.Int)
		s2BigInt := new(big.Int)
		r2.BigInt(r2BigInt)
		s2.BigInt(s2BigInt)

		fmt.Printf("  Signature r (decimal): %s\n", r2BigInt.String())
		fmt.Printf("  Signature s (decimal): %s\n", s2BigInt.String())
		fmt.Printf("  Signature r (hex): 0x%x\n", r2BigInt)
		fmt.Printf("  Signature s (hex): 0x%x\n", s2BigInt)
	}
	fmt.Println()

	// Test Case 3: A larger hash (typical order hash size)
	testHash3 := "3536226799039598760818908146505685573201184860355101576166202238293857277374"
	hash3BigInt, _ := new(big.Int).SetString(testHash3, 10)
	hash3Felt := new(felt.Felt).SetBigInt(hash3BigInt)

	fmt.Printf("Test 3 (Python successful order hash):\n")
	fmt.Printf("  Hash (decimal): %s\n", testHash3)
	fmt.Printf("  Hash (hex): 0x%x\n", hash3BigInt)

	r3, s3, err := curve.SignFelts(hash3Felt, privateKeyFelt)
	if err != nil {
		fmt.Printf("  ERROR signing: %v\n", err)
	} else {
		r3BigInt := new(big.Int)
		s3BigInt := new(big.Int)
		r3.BigInt(r3BigInt)
		s3.BigInt(s3BigInt)

		fmt.Printf("  Signature r (decimal): %s\n", r3BigInt.String())
		fmt.Printf("  Signature s (decimal): %s\n", s3BigInt.String())
		fmt.Printf("  Signature r (hex): 0x%x\n", r3BigInt)
		fmt.Printf("  Signature s (hex): 0x%x\n", s3BigInt)
	}
	fmt.Println()

	fmt.Println("=== Instructions for Python Comparison ===")
	fmt.Println("In Python, use:")
	fmt.Println()
	fmt.Println("from x10.utils.starkex import sign")
	fmt.Println("private_key = <your_private_key_as_int>")
	fmt.Println()
	fmt.Println("# Test 1")
	fmt.Println("r1, s1 = sign(private_key, 1234567890)")
	fmt.Println("print(f'Test 1: r={r1}, s={s1}')")
	fmt.Println("print(f'Test 1 (hex): r=0x{r1:x}, s=0x{s1:x}')")
	fmt.Println()
	fmt.Println("# Test 2")
	fmt.Println("r2, s2 = sign(private_key, 9876543210)")
	fmt.Println("print(f'Test 2: r={r2}, s={s2}')")
	fmt.Println("print(f'Test 2 (hex): r=0x{r2:x}, s=0x{s2:x}')")
	fmt.Println()
	fmt.Println("# Test 3")
	fmt.Println("r3, s3 = sign(private_key, 3536226799039598760818908146505685573201184860355101576166202238293857277374)")
	fmt.Println("print(f'Test 3: r={r3}, s={s3}')")
	fmt.Println("print(f'Test 3 (hex): r=0x{r3:x}, s=0x{s3:x}')")
}
