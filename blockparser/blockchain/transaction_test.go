package blockchain

import "testing"

func TestConvertHexToString(t *testing.T) {
	tests := []struct {
		name string
		hex  string
		want int64
	}{
		{
			name: "valid hex",
			hex:  "0x1",
			want: 1,
		},
		{
			name: "invalid hex",
			hex:  "0xg",
			want: 0,
		},
		{
			name: "empty hex",
			hex:  "",
			want: 0,
		},
		{
			name: "invalid (nil) string",
			hex:  "nil",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertHexToInt(tt.hex); got != tt.want {
				t.Errorf("ConvertHexToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
