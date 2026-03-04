package hasura

import (
	"testing"

	"github.com/ettle/strcase"
)

// TestToSnakeCase verifies that strcase.ToSnake produces the same output as the
// previous hand-rolled ToSnakeCase implementation it replaced. The test guards
// against silent behaviour changes when the dependency is updated.
func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "test 1",
			str:  "L1HandlerCount",
			want: "l1_handler_count",
		}, {
			name: "test 2",
			str:  "Count",
			want: "count",
		}, {
			name: "test 3",
			str:  "ParsedCalldata",
			want: "parsed_calldata",
		}, {
			name: "test 4",
			str:  "L2Block",
			want: "l2_block",
		}, {
			name: "test 5",
			str:  "UserID",
			want: "user_id",
		}, {
			name: "test 6",
			str:  "TxHash",
			want: "tx_hash",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strcase.ToSnake(tt.str); got != tt.want {
				t.Errorf("strcase.ToSnake() = %v, want %v", got, tt.want)
			}
		})
	}
}
