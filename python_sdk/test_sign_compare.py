#!/usr/bin/env python3
"""Compare signing between Python and Go"""

try:
    from fast_stark_crypto import sign
    from vendor.starkware.crypto.signature import generate_k_rfc6979
    print("✓ Successfully imported fast_stark_crypto and vendor")
except ImportError as e:
    print(f"ERROR: Failed to import required modules: {e}")
    exit(1)

# Private key - REPLACE WITH YOUR VALUE FROM .env
# Get it from: X10_PRIVATE_KEY=0x...
private_key = int("0x57b78c7a6db039233ad693bfb625d974de43dda3928c6ab83e58879f02efb83", 16)

print("=== Python Starknet Signing Test (with RFC6979) ===")
print()

# Test 1
print("Test 1:")
msg_hash_1 = 1234567890
print(f"  Hash (decimal): {msg_hash_1}")
print(f"  Hash (hex): {hex(msg_hash_1)}")

k1 = generate_k_rfc6979(msg_hash=msg_hash_1, priv_key=private_key)
r1, s1 = sign(private_key=private_key, msg_hash=msg_hash_1, k=k1)
print(f"  Signature r (decimal): {r1}")
print(f"  Signature s (decimal): {s1}")
print(f"  Signature r (hex): {hex(r1)}")
print(f"  Signature s (hex): {hex(s1)}")
print()

# Test 2
print("Test 2:")
msg_hash_2 = 9876543210
print(f"  Hash (decimal): {msg_hash_2}")
print(f"  Hash (hex): {hex(msg_hash_2)}")

k2 = generate_k_rfc6979(msg_hash=msg_hash_2, priv_key=private_key)
r2, s2 = sign(private_key=private_key, msg_hash=msg_hash_2, k=k2)
print(f"  Signature r (decimal): {r2}")
print(f"  Signature s (decimal): {s2}")
print(f"  Signature r (hex): {hex(r2)}")
print(f"  Signature s (hex): {hex(s2)}")
print()

# Test 3
print("Test 3 (Python successful order hash):")
msg_hash_3 = 3536226799039598760818908146505685573201184860355101576166202238293857277374
print(f"  Hash (decimal): {msg_hash_3}")
print(f"  Hash (hex): {hex(msg_hash_3)}")

k3 = generate_k_rfc6979(msg_hash=msg_hash_3, priv_key=private_key)
r3, s3 = sign(private_key=private_key, msg_hash=msg_hash_3, k=k3)
print(f"  Signature r (decimal): {r3}")
print(f"  Signature s (decimal): {s3}")
print(f"  Signature r (hex): {hex(r3)}")
print(f"  Signature s (hex): {hex(s3)}")
print()

print("=== COMPARISON ===")
print("Now compare these signatures with the Go output from:")
print("  go run examples/test_signing.go")
print()
print("If the signatures MATCH: signing is compatible ✓")
print("If the signatures DIFFER: Go needs RFC6979 implementation")

