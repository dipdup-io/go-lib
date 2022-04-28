package tools

import (
	"testing"
)

func TestIsContract(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "KT1HBy1L43tiLe5MVJZ5RoxGy53Kx8kMgyoU",
			address: "KT1HBy1L43tiLe5MVJZ5RoxGy53Kx8kMgyoU",
			want:    true,
		}, {
			name:    "tz1dMH7tW7RhdvVMR4wKVFF1Ke8m8ZDvrTTE",
			address: "tz1dMH7tW7RhdvVMR4wKVFF1Ke8m8ZDvrTTE",
			want:    false,
		}, {
			name:    "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			address: "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			want:    true,
		}, {
			name:    "expru2dKqDfZG8hu4wNGkiyunvq2hdSKuVYtcKta7BWP6Q18oNxKjS",
			address: "expru2dKqDfZG8hu4wNGkiyunvq2hdSKuVYtcKta7BWP6Q18oNxKjS",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContract(tt.address); got != tt.want {
				t.Errorf("IsContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBigMapKeyHash(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			str:  "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			want: false,
		}, {
			name: "exprtqoNj2hRg8PsPMaXLcy3dXjMM3B7nHKrRNqpfjbYpMbULbRj8k",
			str:  "exprtqoNj2hRg8PsPMaXLcy3dXjMM3B7nHKrRNqpfjbYpMbULbRj8k",
			want: true,
		}, {
			name: "expru2dKqDfZG8hu4wNGkiyunvq2hdSKuVYtcKta7BWP6Q18oNxKjS",
			str:  "expru2dKqDfZG8hu4wNGkiyunvq2hdSKuVYtcKta7BWP6Q18oNxKjS",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBigMapKeyHash(tt.str); got != tt.want {
				t.Errorf("IsBigMapKeyHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOperationHash(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			str:  "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			want: false,
		}, {
			name: "exprtqoNj2hRg8PsPMaXLcy3dXjMM3B7nHKrRNqpfjbYpMbULbRj8k",
			str:  "exprtqoNj2hRg8PsPMaXLcy3dXjMM3B7nHKrRNqpfjbYpMbULbRj8k",
			want: false,
		}, {
			name: "opDqhqYmqgmXTxcEcDXbJMWBThZkaQCovwV8BC3gwthEWYdPCWD",
			str:  "opDqhqYmqgmXTxcEcDXbJMWBThZkaQCovwV8BC3gwthEWYdPCWD",
			want: true,
		}, {
			name: "opRRiHEQacoet5rq7jgcd33K66bkj5qCdThxGnCQwyZtdFjZ8ph",
			str:  "opRRiHEQacoet5rq7jgcd33K66bkj5qCdThxGnCQwyZtdFjZ8ph",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOperationHash(tt.str); got != tt.want {
				t.Errorf("IsOperationHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAddress(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			str:  "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			want: true,
		}, {
			name: "exprtqoNj2hRg8PsPMaXLcy3dXjMM3B7nHKrRNqpfjbYpMbULbRj8k",
			str:  "exprtqoNj2hRg8PsPMaXLcy3dXjMM3B7nHKrRNqpfjbYpMbULbRj8k",
			want: false,
		}, {
			name: "opDqhqYmqgmXTxcEcDXbJMWBThZkaQCovwV8BC3gwthEWYdPCWD",
			str:  "opDqhqYmqgmXTxcEcDXbJMWBThZkaQCovwV8BC3gwthEWYdPCWD",
			want: false,
		}, {
			name: "opRRiHEQacoet5rq7jgcd33K66bkj5qCdThxGnCQwyZtdFjZ8ph",
			str:  "opRRiHEQacoet5rq7jgcd33K66bkj5qCdThxGnCQwyZtdFjZ8ph",
			want: false,
		}, {
			name: "tz1PUnJ3m435ZK4RTqhTEiSYF22YAUx5rEU1",
			str:  "tz1PUnJ3m435ZK4RTqhTEiSYF22YAUx5rEU1",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAddress(tt.str); got != tt.want {
				t.Errorf("IsAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
