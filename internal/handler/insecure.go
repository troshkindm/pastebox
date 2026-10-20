package handler

import (
	"crypto/tls"
	"net/http"
)

// InsecureClient возвращает HTTP-клиент без проверки TLS-сертификатов.
func InsecureClient() *http.Client {
	return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
}
