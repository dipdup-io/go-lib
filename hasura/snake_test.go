package hasura

import "testing"

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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.str); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
