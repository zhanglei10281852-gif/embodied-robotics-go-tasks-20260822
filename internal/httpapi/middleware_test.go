package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTenantHeader(t *testing.T) {
	h := TenantHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	r := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, r)
	if rw.Code != 400 {
		t.Fatal(rw.Code)
	}
	r.Header.Set("X-Tenant-ID", "t")
	rw = httptest.NewRecorder()
	h.ServeHTTP(rw, r)
	if rw.Code != 204 {
		t.Fatal(rw.Code)
	}
}
