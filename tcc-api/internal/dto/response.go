package dto

type ResponseDTO struct {
	Success       bool   `json:"success"`
	Iterations    int    `json:"iterations"`
	ResultingCode string `json:"resulting_code"`
}
