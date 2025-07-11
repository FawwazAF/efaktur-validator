package api

import (
	"net/http"
	"time"
)

type Repository struct {
	HTTPClient *http.Client
	mockURL    string
}

func New() *Repository {
	return &Repository{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		mockURL:    "http://svc.efaktur.pajak.go.id/validasi/faktur/approvalCode/527d5baf11452b2a424b8b899e549f99426cc89fe072d84cac822e58bdf8bb5",
	}
}
