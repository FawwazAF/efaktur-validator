package model

type DJPEfaktur struct {
	SellerTaxID       interface{} `xml:"npwpPenjual"`
	SellerName        interface{} `xml:"namaPenjual"`
	BuyerTaxID        interface{} `xml:"npwpLawanTransaksi"`
	BuyerName         interface{} `xml:"namaLawanTransaksi"`
	EInvoiceNumber    interface{} `xml:"nomorFaktur"`
	EInvoiceDate      interface{} `xml:"tanggalFaktur"`
	TotalTaxBaseValue interface{} `xml:"jumlahPpn"`
	TotalVATAmount    interface{} `xml:"jumlahDpp"`
}
