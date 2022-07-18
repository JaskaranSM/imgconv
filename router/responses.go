package router

type ErrorResponse struct {
	Error string `json:"error"`
}

type ConversionStatusResponse struct {
	Status   string `json:"status"`
	Filename string `json:"filename"`
	Format   string `json:"format"`
	Id       string `json:"id"`
	Error    string `json:"error"`
}
