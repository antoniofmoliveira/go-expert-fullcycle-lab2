package shared

import (
	"testing"
)

func TestValidateStateShort(t *testing.T) {
	type args struct {
		state string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "validate state short",
			args: args{
				state: "MG",
			},
			want: true,
		},
		{
			name: "validate state short error",
			args: args{
				state: "MM",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStateShort(tt.args.state); got != tt.want {
				t.Errorf("ValidateStateShort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateStateLong(t *testing.T) {
	type args struct {
		state string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "validate state long",
			args: args{
				state: "Minas Gerais",
			},
			want: true,
		},
		{
			name: "validate state long error",
			args: args{
				state: "Minas",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStateLong(tt.args.state); got != tt.want {
				t.Errorf("ValidateStateLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateRegiao(t *testing.T) {
	type args struct {
		regiao string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "validate regiao",
			args: args{
				regiao: "Sul",
			},
			want: true,
		},
		{
			name: "validate regiao error",
			args: args{
				regiao: "Sudoeste",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRegiao(tt.args.regiao); got != tt.want {
				t.Errorf("ValidateRegiao() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCep(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "validate cep",
			args: args{
				cep: "39408078",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "validate cep with dash",
			args: args{
				cep: "39408-078",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "validate cep error",
			args: args{
				cep: "3940807",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "validate cep error",
			args: args{
				cep: "3940807-8",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateCep(tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCepWithDash(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "validate cep with dash",
			args: args{
				cep: "39408-078",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "validate cep with dash error",
			args: args{
				cep: "39408078",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateCepWithDash(tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCepWithDash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCepWithDash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCepWithoutDash(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "validate cep without dash",
			args: args{
				cep: "39408078",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "validate cep without dash error",
			args: args{
				cep: "39408-078",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateCepWithoutDash(tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCepWithoutDash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCepWithoutDash() = %v, want %v", got, tt.want)
			}
		})
	}
}
