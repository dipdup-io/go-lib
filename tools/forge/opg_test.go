package forge

import (
	"encoding/hex"
	"testing"

	"github.com/dipdup-net/go-lib/node"
	"github.com/stretchr/testify/assert"
)

func TestOPG(t *testing.T) {
	tests := []struct {
		name       string
		branch     string
		operations []node.Operation
		want       string
		wantErr    bool
	}{
		{
			name:   "test 1",
			branch: "BLRYV1w71DtjyDU27e2XWZ2KyfcGupo985qvphm7PSCNZXk6SHL",
			operations: []node.Operation{
				{
					Kind: node.KindTransaction,
					Body: node.Transaction{
						Amount:       "1000",
						Counter:      "393218",
						Destination:  "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
						Fee:          "351",
						GasLimit:     "1521",
						Source:       "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
						StorageLimit: "100",
					},
				},
				{
					Kind: node.KindTransaction,
					Body: node.Transaction{
						Amount:       "1000",
						Counter:      "393219",
						Destination:  "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
						Fee:          "351",
						GasLimit:     "1521",
						Source:       "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
						StorageLimit: "100",
					},
				},
				{
					Kind: node.KindTransaction,
					Body: node.Transaction{
						Amount:       "1000",
						Counter:      "393220",
						Destination:  "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
						Fee:          "351",
						GasLimit:     "1521",
						Source:       "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
						StorageLimit: "100",
					},
				},
			},
			want: "5db044c1a354b21ef464a61febad3c4efc910588e8f9400d82a64626966af7506c00e8b36c80efb51ec85a14562426049aa182a3ce38df02828018f10b64e8070000e8b36c80efb51ec85a14562426049aa182a3ce38006c00e8b36c80efb51ec85a14562426049aa182a3ce38df02838018f10b64e8070000e8b36c80efb51ec85a14562426049aa182a3ce38006c00e8b36c80efb51ec85a14562426049aa182a3ce38df02848018f10b64e8070000e8b36c80efb51ec85a14562426049aa182a3ce3800",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OPG(tt.branch, tt.operations...)
			if (err != nil) != tt.wantErr {
				t.Errorf("OPG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotHex := hex.EncodeToString(got)
			assert.Equal(t, tt.want, gotHex)
		})
	}
}
