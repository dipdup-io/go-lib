# tools

Core Tezos utilities: Michelson AST, encoding, forging, cryptography, contract interface detection, and code generation helpers.

```bash
go get github.com/dipdup-io/go-lib/tools
```

## Sub-packages

| Package | Import path | Description |
|---------|-------------|-------------|
| `ast` | `.../tools/ast` | Michelson AST — parse, fold, convert, encode |
| `base` | `.../tools/base` | Base node type shared across AST |
| `consts` | `.../tools/consts` | Tezos primitive constants (opcodes, annotations) |
| `contract` | `.../tools/contract` | Contract parser and FA interface detection |
| `crypto` | `.../tools/crypto` | Key generation, signing, signature verification |
| `encoding` | `.../tools/encoding` | Base58Check for addresses, keys, hashes |
| `forge` | `.../tools/forge` | Binary forge/unforge for operations and Michelson |
| `formatter` | `.../tools/formatter` | Michelson source code formatter |
| `tezerrors` | `.../tools/tezerrors` | Typed Tezos protocol error structures |
| `tezgen` | `.../tools/tezgen` | Runtime types for `tezgen`-generated bindings |
| `translator` | `.../tools/translator` | Converts Michelson values between representations |
| `types` | `.../tools/types` | Shared type definitions |

---

## `ast` — Michelson AST

Parse Michelson scripts and work with their typed tree representation.

```go
import "github.com/dipdup-io/go-lib/tools/ast"

// Parse a contract script (parameter/storage/code sections)
script, err := ast.NewScript(rawScriptJSON)

// Unfold a storage value into the AST
if err := script.Storage.FromBigMap(storageJSON); err != nil {
    panic(err)
}

// Convert to the "Miguel" API format (used by TzKT)
miguel, err := script.Storage.ToMiguel()

// Convert to JSON schema
schema, err := script.Parameter.ToJSONSchema()
```

### Contract interface detection

Detect whether a contract implements an FA standard:

```go
iface, err := script.Parameter.FindInterface()
switch iface {
case ast.FA1Interface:
    // FA1 token
case ast.FA12Interface:
    // FA1.2 token
case ast.FA2Interface:
    // FA2 multi-token
}
```

---

## `crypto` — Keys and signatures

```go
import "github.com/dipdup-io/go-lib/tools/crypto"

// Generate a new key (ed25519)
key, err := crypto.NewKey()

// Load from base58-encoded private key
key, err := crypto.NewKeyFromBase58("edsk...")

// Sign arbitrary bytes
sig, err := key.Sign(data)

// Verify a signature
valid, err := key.Verify(data, sig)

// Get the tz1/tz2/tz3 address
addr := key.Address()
```

Supported curves: `ed25519` (tz1), `secp256k1` (tz2), `p256` (tz3).

---

## `encoding` — Base58Check encoding

```go
import "github.com/dipdup-io/go-lib/tools/encoding"

// Encode a public key hash to a tz1 address
addr, err := encoding.EncodeAddress(pubKeyHash, encoding.PrefixED25519PublicKeyHash)

// Decode a tz1 address back to bytes
bytes, err := encoding.DecodeAddress("tz1...")

// Encode an operation hash
hash, err := encoding.EncodeOperationHash(opBytes)
```

---

## `forge` — Binary encoding

Encode and decode Tezos binary formats.

```go
import "github.com/dipdup-io/go-lib/tools/forge"

// Forge a natural number (variable-length encoding)
encoded := forge.ForgeNat(12345)

// Forge a signed integer
encoded = forge.ForgeInt(-42)

// Forge a boolean
encoded = forge.ForgeBool(true)

// Unforge (decode) an operation from hex
operation, err := forge.UnforgeOperation("deadbeef...")

// Forge a complete operation
raw, err := forge.ForgeOperation(operation)
```

---

## `contract` — Contract parsing

```go
import "github.com/dipdup-io/go-lib/tools/contract"

// Parse contract script and detect tags/interfaces
tags, err := contract.GetTags(script)

// Parse a specific entrypoint's parameter type
paramType, err := contract.ParseParameter(script, "transfer")
```

### FA interface detection via tags

```go
if tags.Has(contract.FA2Tag) {
    // contract implements FA2
}
```

---

## `tezerrors` — Protocol errors

Typed wrappers for errors returned by the Tezos protocol:

```go
import "github.com/dipdup-io/go-lib/tools/tezerrors"

errs, err := tezerrors.ParseErrors(rawJSON)
for _, e := range errs {
    if tezerrors.HasError(errs, "proto.016.insufficient_balance") {
        // handle specific error
    }
}
```

---

## `tezgen` — Generated contract binding types

Runtime types used by contracts generated with the `tezgen` tool:

```go
import "github.com/dipdup-io/go-lib/tools/tezgen"

// Types available: Address, Bytes, Int, Contract,
// SaplingState, SaplingTransaction, Timestamp, Unit
var addr tezgen.Address
if err := addr.UnmarshalJSON(data); err != nil {
    panic(err)
}
```
