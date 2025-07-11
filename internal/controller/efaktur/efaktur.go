package efaktur

import (
	"bufio"
	"context"
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/efaktur-validator/internal/model"
	"github.com/unidoc/unipdf/v4/extractor"
	pdfModel "github.com/unidoc/unipdf/v4/model"
)

func (ef *efakturController) ValidateEfaktur(ctx context.Context, pdfPath string) (model.EfakturValidationResult, error) {
	// parse PDF request
	parsedPDF, err := ef.ParseEfakturPDF(ctx, pdfPath)
	if err != nil {
		return model.EfakturValidationResult{}, err
	}

	url := parsedPDF.QRUrl
	djpInvoices, err := ef.djp.GetInvoicesFromDJP(ctx, url)
	if err != nil {
		return model.EfakturValidationResult{}, err
	}

	return ef.CompareRequestEfakturWithDKP(ctx, parsedPDF, djpInvoices), nil
}

// TODO: changes new lib for parsing PDF
func (ef *efakturController) ParseEfakturPDF(ctx context.Context, pdfPath string) (model.EfakturPDF, error) {
	filePDF, err := os.Open(pdfPath)
	if err != nil {
		return model.EfakturPDF{}, err
	}
	defer filePDF.Close()

	pdfReader, err := pdfModel.NewPdfReader(filePDF)
	if err != nil {
		return model.EfakturPDF{}, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return model.EfakturPDF{}, err
	}

	if numPages < 1 {
		return model.EfakturPDF{}, errors.New("got empty pdf")
	}

	page, err := pdfReader.GetPage(1) // assuming efaktur only have 1 page.
	if err != nil {
		return model.EfakturPDF{}, err
	}

	ex, err := extractor.New(page)
	if err != nil {
		return model.EfakturPDF{}, err
	}

	text, err := ex.ExtractText()
	if err != nil {
		return model.EfakturPDF{}, err
	}

	result := model.EfakturPDF{}
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "NPWP Penjual:"):
			result.SellerTaxID = strings.TrimSpace(strings.TrimPrefix(line, "NPWP Penjual:"))
		case strings.HasPrefix(line, "Nama Penjual:"):
			result.SellerName = strings.TrimSpace(strings.TrimPrefix(line, "Nama Penjual:"))
		case strings.HasPrefix(line, "NPWP Pembeli:"):
			result.BuyerTaxID = strings.TrimSpace(strings.TrimPrefix(line, "NPWP Pembeli:"))
		case strings.HasPrefix(line, "Nama Pembeli:"):
			result.BuyerName = strings.TrimSpace(strings.TrimPrefix(line, "Nama Pembeli:"))
		case strings.HasPrefix(line, "Nomor Faktur:"):
			result.EInvoiceNumber = strings.TrimSpace(strings.TrimPrefix(line, "Nomor Faktur:"))
		case strings.HasPrefix(line, "Tanggal Faktur:"):
			result.EInvoiceDate = strings.TrimSpace(strings.TrimPrefix(line, "Tanggal Faktur:"))
		case strings.HasPrefix(line, "DPP:"):
			result.TotalTaxBaseValue = strings.TrimSpace(strings.TrimPrefix(line, "DPP:"))
		case strings.HasPrefix(line, "PPN:"):
			result.TotalVATAmount = strings.TrimSpace(strings.TrimPrefix(line, "PPN:"))
		case strings.HasPrefix(line, "QR Code:"):
			result.QRUrl = strings.TrimSpace(strings.TrimPrefix(line, "QR Code:"))
		}
	}

	//TODO : validate extracted images
	return result, nil
}

func (ef *efakturController) CompareRequestEfakturWithDKP(ctx context.Context, reqInvoice model.EfakturPDF, djpInvoice model.DJPEfaktur) model.EfakturValidationResult {
	deviations := make([]model.Deviation, 0)
	validated := make(map[string]interface{})

	compareField := func(field string, pdfVal, apiVal any) {
		if pdfVal == nil || apiVal == nil {
			// One is missing
			if pdfVal == nil {
				deviations = append(deviations, model.Deviation{
					Field:         field,
					PDFValue:      nil,
					DJPAPIValue:   apiVal,
					DeviationType: "missing_in_pdf",
				})
			} else {
				deviations = append(deviations, model.Deviation{
					Field:         field,
					PDFValue:      pdfVal,
					DJPAPIValue:   nil,
					DeviationType: "missing_in_api",
				})
			}
			return
		}

		if !reflect.DeepEqual(pdfVal, apiVal) {
			deviations = append(deviations, model.Deviation{
				Field:         field,
				PDFValue:      pdfVal,
				DJPAPIValue:   apiVal,
				DeviationType: "mismatch",
			})
			return
		}

		// Validated match
		validated[field] = apiVal
	}

	compareField("SellerTaxID", reqInvoice.SellerTaxID, djpInvoice.SellerTaxID)
	compareField("SellerName", reqInvoice.SellerName, djpInvoice.SellerName)
	compareField("BuyerTaxID", reqInvoice.BuyerTaxID, djpInvoice.BuyerTaxID)
	compareField("BuyerName", reqInvoice.BuyerName, djpInvoice.BuyerName)
	compareField("EInvoiceNumber", reqInvoice.EInvoiceNumber, djpInvoice.EInvoiceNumber)
	compareField("EInvoiceDate", reqInvoice.EInvoiceDate, djpInvoice.EInvoiceDate)
	compareField("TotalTaxBaseValue", reqInvoice.TotalTaxBaseValue, djpInvoice.TotalTaxBaseValue)
	compareField("TotalVATAmount", reqInvoice.TotalVATAmount, djpInvoice.TotalVATAmount)

	status := "validated_successfully"
	if len(deviations) > 0 {
		status = "validated_with_deviations"
	}

	return model.EfakturValidationResult{
		Status:  status,
		Message: "success",
		ValidationResults: model.ValidationResults{
			Deviations:    deviations,
			ValidatedData: validated,
		},
	}
}
