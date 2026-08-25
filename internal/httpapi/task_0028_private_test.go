package httpapi

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/service"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
)

func TestMissionHandlerStopsOnCancelledRequest(t *testing.T) {
	ctx := context.Background(); db, err := storage.Open(ctx, "file:http-cancel-private?mode=memory&cache=shared"); if err != nil { t.Fatal(err) }; defer db.Close()
	r := repository.New(db); _ = r.CreateTenant(ctx, domain.Tenant{ID:"tenant",Name:"t",Quota:1}); _ = r.Robots().Create(ctx, domain.Robot{ID:"http-r",TenantID:"tenant",Serial:"HTTP",Name:"http",Status:domain.RobotReady,Version:1})
	h := New(service.New(r)); cancelled, cancel := context.WithCancel(ctx); cancel(); req := httptest.NewRequest("POST", "/v1/missions", strings.NewReader("{\"tenantID\":\"tenant\",\"robotID\":\"http-r\",\"priority\":1,\"policyVersion\":1}")).WithContext(cancelled); rw := httptest.NewRecorder(); h.createMission(rw, req)
	if rw.Code < 400 { t.Fatalf("cancelled request created mission: %d", rw.Code) }
}
