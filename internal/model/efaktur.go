package model

type EfakturPDF struct {
	SellerTaxID       interface{}
	SellerName        interface{}
	BuyerTaxID        interface{}
	BuyerName         interface{}
	EInvoiceNumber    interface{}
	EInvoiceDate      interface{}
	TotalTaxBaseValue interface{}
	TotalVATAmount    interface{}
	QRUrl             string
}

type EfakturValidationResult struct {
	Status            string            `json:"status"`
	Message           string            `json:"message"`
	ValidationResults ValidationResults `json:"validation_results"`
}

type ValidationResults struct {
	Deviations    []Deviation            `json:"deviations"`
	ValidatedData map[string]interface{} `json:"validated_data"`
}

type Deviation struct {
	Field         string      `json:"field"`
	PDFValue      interface{} `json:"pdf_value"`
	DJPAPIValue   interface{} `json:"djp_api_value"`
	DeviationType string      `json:"deviation_type"`
}
