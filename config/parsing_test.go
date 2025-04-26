package config

import (
	"os"
	"testing"

	"golang.org/x/text/transform"
)

func Test_expandTransformer_Transform(t *testing.T) {
	tests := []struct {
		name    string
		dst     []byte
		src     []byte
		envs    map[string]string
		want    int
		want1   int
		wantErr bool
		atEOF   bool
	}{
		{
			name:    "test 1",
			dst:     make([]byte, 256),
			src:     []byte(`$tring without {env} `),
			atEOF:   false,
			envs:    map[string]string{},
			want:    21,
			want1:   21,
			wantErr: false,
		}, {
			name:  "test 2",
			dst:   make([]byte, 256),
			src:   []byte(`string with ${TEST_ENV} `),
			atEOF: false,
			envs: map[string]string{
				"TEST_ENV": "test_value",
			},
			want:    23,
			want1:   24,
			wantErr: false,
		}, {
			name:    "test 3",
			dst:     make([]byte, 256),
			src:     []byte(`string with ${TEST_ENV:-defalt_value} `),
			atEOF:   false,
			envs:    map[string]string{},
			want:    25,
			want1:   38,
			wantErr: false,
		}, {
			name:  "test 4",
			dst:   make([]byte, 256),
			src:   []byte(`string with ${TEST_ENV:-defalt_value} ${TEST_ENV2:-defalt_value2}`),
			atEOF: false,
			envs: map[string]string{
				"TEST_ENV": "test_value",
			},
			want:    36,
			want1:   65,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				for key := range tt.envs {
					os.Unsetenv(key)
				}
			}()
			for key, value := range tt.envs {
				t.Setenv(key, value)
			}

			transformer := &expandTransformer{
				NopResetter: transform.NopResetter{},
			}
			got, got1, err := transformer.Transform(tt.dst, tt.src, tt.atEOF)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandTransformer.Transform() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("expandTransformer.Transform() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("expandTransformer.Transform() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
