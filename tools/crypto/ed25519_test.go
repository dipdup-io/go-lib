package crypto

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/dipdup-net/go-lib/tools/encoding"
)

func TestEd25519_Sign(t *testing.T) {
	type args struct {
		data       string
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				data:       "03f265006be7d2f678d4ca8cae9f1877b40a76a21f8b05e73254043af01f3638786c00739ab9281b15479d756572217dc8bb944f2b02a7880ef0e42ac350f403c0843d00008384a29947e81770a1586cb268552a25ee81fb7a00",
				privateKey: "edsk2tbfnJFHrV7R8XxA8fKGFdsPotvkPLDqTa12XMMQ9v8yZ5A7zy",
			},
			want: "edsigtjJQ2wppYpkR16kk4QU6mH758tUoSahfPcykQxpVaU1wJ5HfBy1tANeu9BL7kWYdrGUfqWubYtiU6HNNrV8mtUWNScBWAu",
		}, {
			name: "test 2",
			args: args{
				data:       "03f265006be7d2f678d4ca8cae9f1877b40a76a21f8b05e73254043af01f3638786c0027b104098f9b1b8e84261d03f4f788841bc44eb7880ec2d02ac350f403c0843d00008384a29947e81770a1586cb268552a25ee81fb7a00",
				privateKey: "edsk4Mz9VrUKZBXBJ3EkguGxpmwHZSiefz8SibjKko5qjpwGvLDbum",
			},
			want: "edsigtqmbEzia3HvtZ4i5WsZ9ci1KNe7pseEJA3TUr6hyPpkp58KZhqhrhG2CRcbcftr1FVCEGg6KsVUVCyc9hxNL6HKPs7Rqug",
		}, {
			name: "test 3",
			args: args{
				data:       "03f265006be7d2f678d4ca8cae9f1877b40a76a21f8b05e73254043af01f3638786c00fe2ce0cccc0214af521ad60c140c5589b4039247880ef0e42ac350f403c0843d00008384a29947e81770a1586cb268552a25ee81fb7a00",
				privateKey: "edsk35mfZXZJiYUxqcmsK5K6ggg3owD2dpbRgFHp4zZzmrPy9RBdj8",
			},
			want: "edsigtszSUdsrNe3jM8YAxXHVSxfS7L5NzippwjtfxBAM86FcrXzZnURtvEqnvUEZd9sXbN2fUQn9ZnHe3L7wfAkABJYgDXNNtY",
		}, {
			name: "test 4",
			args: args{
				data:       "03f265006be7d2f678d4ca8cae9f1877b40a76a21f8b05e73254043af01f3638786c009d923eb28e4255e0171a63304363b904f78b36bd880eefe42ac350f403c0843d00008384a29947e81770a1586cb268552a25ee81fb7a00",
				privateKey: "edsk4CWL88P8PxhXojtxUxq8u9NpPFzSjCCRJJUWGi2bD7C6juiYwF",
			},
			want: "edsigtbuLCYQ79DMukbVHqH5i6Svsh5Agn8gg7DT5fBAGGr4qsTEHTveU4uDnKmUK4zDPBXU9bTR45RQWV2emNRQJpaiLhYEjYh",
		}, {
			name: "test 5",
			args: args{
				data:       "03f265006be7d2f678d4ca8cae9f1877b40a76a21f8b05e73254043af01f3638786c004050e870eeeec5897df25de422b8da1b255dde14880eefe42ac350f403c0843d00008384a29947e81770a1586cb268552a25ee81fb7a00",
				privateKey: "edsk2kL1z6s7Lo6z2YfZQfX2RzoLktoNezKuBQ4vgvAPFBM2J4BcfM",
			},
			want: "edsigtY11avmRNuycJU8MbkJczKeWy6RSpVwGzrmrZUw15JXPqdD98RCtCM9Q194B19YrBGBDV2VzbxDvQBwK1p3cXhUanQUf3v",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curve := NewEd25519()
			data, err := hex.DecodeString(tt.args.data)
			if err != nil {
				t.Errorf("hex.DecodeString(%s) error = %v", tt.args.data, err)
				return
			}

			sk, err := encoding.DecodeBase58(tt.args.privateKey)
			if err != nil {
				t.Errorf("encoding.DecodeBase58(%s) error = %v", tt.args.privateKey, err)
				return
			}

			pk, err := curve.GetPublicKey(sk)
			if err != nil {
				t.Errorf("curve.GetPublicKey(%v) error = %v", sk, err)
				return
			}

			got, err := curve.Sign(data, append(sk, pk...))
			if (err != nil) != tt.wantErr {
				t.Errorf("Ed25519.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			signature, err := encoding.DecodeBase58(tt.want)
			if err != nil {
				t.Errorf("encoding.DecodeBase58(%s) error = %v", tt.want, err)
				return
			}

			if !curve.Verify(data, got.bytes, pk) {
				t.Errorf("Ed25519.Verify() data = %v signature = %v", data, signature)
				return
			}

			if !reflect.DeepEqual(got.bytes, signature) {
				t.Errorf("Ed25519.Sign() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
