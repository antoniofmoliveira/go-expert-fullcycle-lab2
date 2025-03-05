package dto

import (
	"reflect"
	"testing"
)

func TestNewViacep(t *testing.T) {
	type args struct {
		cep         string
		logradouro  string
		complemento string
		unidade     string
		bairro      string
		localidade  string
		uf          string
		estado      string
		regiao      string
		ibge        string
		gia         string
		ddd         string
		siafi       string
	}
	tests := []struct {
		name    string
		args    args
		want    *Viacep
		wantErr bool
	}{
		{
			name: "new viacep",
			args: args{
				cep:         "3940807",
				logradouro:  "Avenida Herlindo Silveira",
				complemento: "Apto 101",
				unidade:     "Sala 101",
				bairro:      "Centro",
				localidade:  "Montes Claros",
				uf:          "MG",
				estado:      "Minas Gerais",
				regiao:      "Sudeste",
				ibge:        "3143302",
				gia:         "",
				ddd:         "38",
				siafi:       "4865",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new viacep wothout error",
			args: args{
				cep:         "39408-078",
				logradouro:  "Avenida Herlindo Silveira",
				complemento: "Apto 101",
				unidade:     "Sala 101",
				bairro:      "Centro",
				localidade:  "Montes Claros",
				uf:          "MG",
				estado:      "Minas Gerais",
				regiao:      "Sudeste",
				ibge:        "3143302",
				gia:         "",
				ddd:         "38",
				siafi:       "4865",
			},
			want: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewViacep(tt.args.cep, tt.args.logradouro, tt.args.complemento, tt.args.unidade, tt.args.bairro, tt.args.localidade, tt.args.uf, tt.args.estado, tt.args.regiao, tt.args.ibge, tt.args.gia, tt.args.ddd, tt.args.siafi)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewViacep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewViacep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewViacepFromJson(t *testing.T) {
	type args struct {
		jsonString string
	}
	tests := []struct {
		name    string
		args    args
		want    *Viacep
		wantErr bool
	}{
		{
			name: "new viacep from json",
			args: args{
				jsonString: `{
					"cep": "39408-078",
					"logradouro": "Avenida Herlindo Silveira",
					"complemento": "Apto 101",
					"unidade": "Sala 101",
					"bairro": "Centro",
					"localidade": "Montes Claros",
					"uf": "MG",
					"estado": "Minas Gerais",
					"regiao": "Sudeste",
					"ibge": "3143302",
					"gia": "",
					"ddd": "38",
					"siafi": "4865"
				}`,
			},
			want: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: false,
		},
		{
			name: "new viacep from json error",
			args: args{
				jsonString: `{
					"cep": 39408-078,
					"logradouro": "Avenida Herlindo Silveira",
					"complemento": "Apto 101",
					"unidade": "Sala 101",
					"bairro": "Centro",
					"localidade": "Montes Claros",
					"uf": "MG",
					"estado": "Minas Gerais",
					"regiao": "Sudeste",
					"ibge": "3143302",
					"gia": "",
					"ddd": "38",
					"siafi": "4865"
				}`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new viacep from json validation error",
			args: args{
				jsonString: `{
					"cep": "39408078",
					"logradouro": "Avenida Herlindo Silveira",
					"complemento": "Apto 101",
					"unidade": "Sala 101",
					"bairro": "Centro",
					"localidade": "Montes Claros",
					"uf": "MG",
					"estado": "Minas Gerais",
					"regiao": "Sudeste",
					"ibge": "3143302",
					"gia": "",
					"ddd": "38",
					"siafi": "4865"
				}`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewViacepFromJson(tt.args.jsonString)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewViacepFromJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewViacepFromJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestViacep_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       *Viacep
		wantErr bool
	}{
		{
			name:    "validate viacep",
			v:       &Viacep{},
			wantErr: true,
		},
		{
			name: "validate viacep cep error",
			v: &Viacep{
				Cep:         "39408078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
		{
			name: "validate viacep uf error",
			v: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MM",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
		{
			name: "validate viacep estado error",
			v: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
		{
			name: "validate viacep regiao error",
			v: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudoeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
		{
			name: "validate viacep logradouro error",
			v: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
		{
			name: "validate viacep bairro error",
			v: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "",
				Localidade:  "Montes Claros",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
		{
			name: "validate viacep localidade error",
			v: &Viacep{
				Cep:         "39408-078",
				Logradouro:  "Avenida Herlindo Silveira",
				Complemento: "Apto 101",
				Unidade:     "Sala 101",
				Bairro:      "Centro",
				Localidade:  "",
				Uf:          "MG",
				Estado:      "Minas Gerais",
				Regiao:      "Sudeste",
				Ibge:        "3143302",
				Gia:         "",
				Ddd:         "38",
				Siafi:       "4865",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Viacep.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
