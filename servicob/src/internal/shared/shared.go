package shared

import (
	"errors"
	"regexp"

	"golang.org/x/exp/slices"
)

var states_short = []string{"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO"}

var states_long = []string{"Acre", "Alagoas",
	"Amapá", "Amazonas", "Bahia", "Ceará", "Distrito Federal",
	"Espirito Santo", "Goiás", "Maranhão", "Mato Grosso", "Mato Grosso do Sul",
	"Minas Gerais", "Paraí", "Pará", "Paraná", "Pernambuco", "Piaui", "Rio de Janeiro",
	"Rio Grande do Norte", "Rio Grande do Sul", "Rondônia", "Roraima", "Santa Catarina",
	"São Paulo", "Sergipe", "Tocantins"}

var regioes = []string{"Sul", "Sudeste", "Centro-Oeste", "Norte", "Nordeste"}

// ValidateCep checks if the given cep is valid.
//
// A valid cep must have 8 digits, optionally with '-'.
// The '-' character is allowed, but not required.
//
// If the cep is valid, it returns true, nil. Otherwise, it returns false, error.
func ValidateCep(cep string) (bool, error) {
	var cepRegex = regexp.MustCompile(`^\d{5}-?\d{3}$`)
	if !cepRegex.MatchString(cep) {
		return false, errors.New("cep must have 8 digits, optionally with '-'")
	}
	return true, nil
}

// ValidateCepWithDash checks if the given cep is valid, with '-'.
//
// A valid cep must have 8 digits, with '-'.
// The '-' character is required.
//
// If the cep is valid, it returns true, nil. Otherwise, it returns false, error.
func ValidateCepWithDash(cep string) (bool, error) {
	var cepRegex = regexp.MustCompile(`^\d{5}-\d{3}$`)
	if !cepRegex.MatchString(cep) {
		return false, errors.New("cep must have 8 digits, optionally with '-'")
	}
	return true, nil
}

// ValidateCepWithoutDash checks if the given cep is valid, without '-'.
//
// A valid cep must have 8 digits, without '-'.
// The '-' character is not allowed.
//
// If the cep is valid, it returns true, nil. Otherwise, it returns false, error.
func ValidateCepWithoutDash(cep string) (bool, error) {
	var cepRegex = regexp.MustCompile(`^\d{8}$`)
	if !cepRegex.MatchString(cep) {
		return false, errors.New("cep must have 8 digits, optionally with '-'")
	}
	return true, nil
}

// ValidateStateShort checks if the given state abbreviation is valid.
//
// A valid state abbreviation must be one of the recognized Brazilian state codes,
// such as "AC" for Acre or "SP" for São Paulo.
//
// If the state abbreviation is valid, it returns true. Otherwise, it returns false.
func ValidateStateShort(state string) bool {
	return slices.Contains(states_short, state)
}

// ValidateStateLong checks if the given state name is valid.
//
// A valid state name must be one of the recognized Brazilian state names,
// such as "Acre" or "São Paulo".
//
// If the state name is valid, it returns true. Otherwise, it returns false.
func ValidateStateLong(state string) bool {
	return slices.Contains(states_long, state)
}

// ValidateRegiao checks if the given region name is valid.
//
// A valid region name must be one of the recognized Brazilian region names,
// such as "Sul" or "Nordeste".
//
// If the region name is valid, it returns true. Otherwise, it returns false.
func ValidateRegiao(regiao string) bool {
	return slices.Contains(regioes, regiao)
}
