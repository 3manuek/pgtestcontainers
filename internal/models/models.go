package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	GetBoundsQuery  string = `SELECT min(ts), max(ts) FROM devices`
	AggregationWOMV string = `SELECT 
    								time_bucket('20 minute', ts) AS minute,
    								device,
    								max(payload->>'cpu') AS cpu
								FROM devices
								GROUP BY minute, device;`
	AggregationWMV string = `SELECT minute, device, cpu from device_cpu_by_minute;`
)

type TS interface {
	FilterByTS(ctx context.Context, filter *Filters) (*sql.Rows, error)
	GetBounds(ctx context.Context) (*Filters, error)
	AggregateMaxCPUMinuteWOMV(ctx context.Context) (*sql.Rows, error)
	AggregateMaxCPUMinuteWMV(ctx context.Context) (*sql.Rows, error)
}

type Aggs struct {
	Minute time.Time
	Device uuid.UUID
	Cpu    float64
}

type Payload struct {
	Cpu  float64 `json:"cpu"`
	Temp float64 `json:"temp"`
}

type Filters struct {
	StartTime time.Time
	EndTime   time.Time
}
type DevicesRow struct {
	Ts time.Time `db:"ts"`
	// // Device  string    `db:"device"` // Handling tdevice enum, disabled for now
	Device  uuid.UUID `db:"device"`
	Payload *Payload  `db:"payload"`
}

type Devices struct {
	Conn   *sql.DB // db conn
	Bounds *Filters
}

func NewDevices(ctx context.Context, db *sql.DB) (Devices, error) {
	d := Devices{Conn: db}
	// ping the db
	err := db.Ping()
	if err != nil {
		return Devices{}, err
	}
	f, err := d.GetBounds(ctx)
	if err != nil {
		log.Fatalln("Failed to get bounds:", err)
		return Devices{}, err
	}

	return Devices{Conn: db, Bounds: &f}, nil
}

func (d *Devices) GetBounds(ctx context.Context) (Filters, error) {
	f := Filters{}

	SingleRowQuery(ctx, d.Conn, GetBoundsQuery, &f.StartTime, &f.EndTime)

	return f, nil
}

func (d *Devices) AggregateMaxCPUMinuteWOMV(ctx context.Context) ([]Aggs, error) {
	a := []Aggs{}
	Query(ctx, d.Conn, AggregationWOMV, &a)
	return a, nil
}

func (d *Devices) AggregateMaxCPUMinuteWMV(ctx context.Context) ([]Aggs, error) {
	a := []Aggs{}
	Query(ctx, d.Conn, AggregationWMV, &a)
	// SingleRowQuery(ctx, d.Conn, AggregationWOMV, &a.Minute, &a.Device, &a.Cpu)
	return a, nil
}

// TODO:
// To implement

func (d *Devices) FilterByTS(ctx context.Context, filter *Filters) (*sql.Rows, error) {

	return nil, nil
}
