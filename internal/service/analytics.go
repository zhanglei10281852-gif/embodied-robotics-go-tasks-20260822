package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"math"
	"sort"
	"time"
)

type Pose struct {
	X, Y, Z float64
	At      time.Time
}
type MotionStats struct {
	Distance, Duration, AverageSpeed float64
	Samples                          int
}

func (s *Service) MotionStatistics(ctx context.Context, tenant, robot string, before time.Time) (MotionStats, error) {
	page, e := s.ReadTelemetry(ctx, tenant, robot, before, 500)
	if e != nil {
		return MotionStats{}, e
	}
	poses := []Pose{}
	for _, ev := range page.Items {
		p, e := s.DecodeTelemetryPayload(ev)
		if e != nil {
			continue
		}
		x, _ := p["x"].(float64)
		y, _ := p["y"].(float64)
		z, _ := p["z"].(float64)
		poses = append(poses, Pose{X: x, Y: y, Z: z, At: ev.RecordedAt})
	}
	sort.Slice(poses, func(i, j int) bool { return poses[i].At.Before(poses[j].At) })
	var distance float64
	for i := 1; i < len(poses); i++ {
		dx := poses[i].X - poses[i-1].X
		dy := poses[i].Y - poses[i-1].Y
		dz := poses[i].Z - poses[i-1].Z
		distance += math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
	duration := 0.0
	if len(poses) > 1 {
		duration = poses[len(poses)-1].At.Sub(poses[0].At).Seconds()
	}
	speed := 0.0
	if duration > 0 {
		speed = distance / duration
	}
	return MotionStats{Distance: distance, Duration: duration, AverageSpeed: speed, Samples: len(poses)}, nil
}

type Availability struct{ Ready, Busy, Offline, Retired int }

func (s *Service) Availability(ctx context.Context, tenant string) (Availability, error) {
	robots, e := s.Repos.Robots().List(ctx, tenant, 1000)
	if e != nil {
		return Availability{}, e
	}
	out := Availability{}
	for _, r := range robots {
		switch r.Status {
		case domain.RobotReady:
			out.Ready++
		case domain.RobotBusy:
			out.Busy++
		case domain.RobotOffline:
			out.Offline++
		case domain.RobotRetired:
			out.Retired++
		}
	}
	return out, nil
}
func Percent(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) * 100 / float64(total)
}
