package forge

import (
	"encoding/hex"
	"testing"

	"github.com/dipdup-net/go-lib/node"
)

func TestTransaction(t *testing.T) {
	tests := []struct {
		name    string
		txJSON  string
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			txJSON: `{
				"kind": "transaction",
				"source": "tz1NXjqkurAmpKJEF76T58oyNsy3hWK7mk8e",
				"fee": "22100",
				"counter": "377727",
				"gas_limit": "218465",
				"storage_limit": "668",
				"amount": "0",
				"destination": "KT1SkmB19o8nfhRvG9LL7TjDfX2Bm1nCuYoY"
			}`,
			want: "6c001fb7d0a599ddca61b88dc203eeefbac341422cdfd4ac01ff8617e1aa0d9c050001c756189bc655cc487d57e5fefe482449dbe00c390000",
		}, {
			name: "test 2",
			txJSON: ` {
				"kind": "transaction",
				"source": "tz1SJJY253HoEda8PS5vvfHVtyghgK3CTS2z",
				"fee": "2966",
				"counter": "133558",
				"gas_limit": "26271",
				"storage_limit": "0",
				"amount": "0",
				"destination": "KT1XdCkJncWfGvqf1NdbK2HBRTvRcHhJtNx5",
				"parameters": {
				  "entrypoint": "do",
				  "value": [
					{
					  "prim": "RENAME"
					},
					{
					  "prim": "NIL",
					  "args": [
						{
						  "prim": "operation"
						}
					  ]
					},
					{
					  "prim": "PUSH",
					  "args": [
						{
						  "prim": "key_hash"
						},
						{
						  "string": "tz2L2HuhaaSnf6ShEDdhTEAr5jGPWPNwpvcB"
						}
					  ]
					},
					{
					  "prim": "IMPLICIT_ACCOUNT"
					},
					{
					  "prim": "PUSH",
					  "args": [
						{
						  "prim": "mutez"
						},
						{
						  "int": "2"
						}
					  ]
					},
					{
					  "prim": "UNIT"
					},
					{
					  "prim": "TRANSFER_TOKENS"
					},
					{
					  "prim": "CONS"
					},
					{
					  "prim": "DIP",
					  "args": [
						[
						  {
							"prim": "DROP"
						  }
						]
					  ]
					}
				  ]
				}
			}`,
			want: "6c00490dc9520ec45270f240a3cc4f07aec76adc358d9617b693089fcd01000001fcc0bee1480bfca3a80481904cee4099400b1c8d00ff020000004f020000004a0358053d036d0743035d0100000024747a324c324875686161536e663653684544646854454172356a475057504e7770766342031e0743036a0002034f034d031b051f02000000020320",
		}, {
			name: "test 3",
			txJSON: ` {
				"kind": "transaction",
				"source": "tz1XJ1UNechmHKhQo4tvVX6qztnVuQuSFKgd",
				"fee": "1283",
				"counter": "7",
				"gas_limit": "10307",
				"storage_limit": "0",
				"amount": "20000000000",
				"destination": "tz1aWXP237BLwNHJcCD4b3DutCevhqq2T1Z9"
			}`,
			want: "6c007fd82c06cf5a203f18faaf562447ed1efcc6c010830a07c350008090dfc04a0000a31e81ac3425310e3274a4698a793b2839dc0afa00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var operation node.Operation
			if err := json.UnmarshalFromString(tt.txJSON, &operation); err != nil {
				t.Errorf("UnmarshalFromString() error = %v", err)
				return
			}
			transaction, err := node.NewTypedOperation[node.Transaction](operation)
			if err != nil {
				t.Errorf("operation.Transaction() error = %v", err)
				return
			}
			got, err := Transaction(transaction)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotStr := hex.EncodeToString(got)
			if tt.want != gotStr {
				t.Errorf("Transaction() = %s, want %s", gotStr, tt.want)
			}
		})
	}
}
