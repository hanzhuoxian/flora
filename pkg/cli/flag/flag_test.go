package flag

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestWordSepNormalizeFunc(t *testing.T) {
	tests := []struct {
		name    string
		want    pflag.NormalizedName
		wantErr bool
	}{
		{
			name:    "Test case 1",
			want:    pflag.NormalizedName("Test case 1"),
			wantErr: false,
		},
		{
			name:    "Test case 2",
			want:    pflag.NormalizedName("Test case 2"),
			wantErr: false,
		},
		{
			name:    "Test_case_3",
			want:    pflag.NormalizedName("Test-case-3"),
			wantErr: false,
		},
		{
			name:    "Test_case_4",
			want:    pflag.NormalizedName("Test-case-4"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordSepNormalizeFunc(pflag.NewFlagSet("test", pflag.ExitOnError), tt.name)

			if got != tt.want {
				t.Errorf("WordSepNormalizeFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
