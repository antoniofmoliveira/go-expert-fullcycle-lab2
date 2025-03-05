package dto

type TempResponse struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
}
