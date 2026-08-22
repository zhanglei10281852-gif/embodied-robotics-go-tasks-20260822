package httpapi

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/service"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthAndRobotHTTP(t *testing.T) {
	db, e := storage.Open(context.Background(), "file:http-test?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := repository.New(db)
	if e = r.CreateTenant(context.Background(), domain.Tenant{ID: "t", Name: "T", Quota: 2}); e != nil {
		t.Fatal(e)
	}
	h := New(service.New(r))
	ts := httptest.NewServer(h.Routes())
	defer ts.Close()
	resp, e := http.Get(ts.URL + "/healthz")
	if e != nil || resp.StatusCode != 200 {
		t.Fatalf("health: %v %v", e, resp.StatusCode)
	}
	body := `{"tenantID":"t","serial":"S","name":"R"}`
	resp, e = http.Post(ts.URL+"/v1/robots", "application/json", strings.NewReader(body))
	if e != nil || resp.StatusCode != 201 {
		t.Fatalf("robot: %v %v", e, resp.StatusCode)
	}
}
