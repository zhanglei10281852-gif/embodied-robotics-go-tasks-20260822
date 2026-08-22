package audit

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestAuditChain(t *testing.T) {
	now := time.Now()
	c := Chain{}
	if e := c.Append(domain.AuditEvent{ID: "1", TenantID: "t", OccurredAt: now}); e != nil {
		t.Fatal(e)
	}
	if e := c.Append(domain.AuditEvent{ID: "2", TenantID: "t", OccurredAt: now.Add(time.Second)}); e != nil {
		t.Fatal(e)
	}
	if c.Digest() == "" || len(c.Since(now)) != 2 {
		t.Fatal("chain")
	}
	if e := c.Append(domain.AuditEvent{ID: "0", TenantID: "t", OccurredAt: now.Add(-time.Second)}); e == nil {
		t.Fatal("order")
	}
}
