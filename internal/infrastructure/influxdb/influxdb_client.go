package influxdb_client

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
)

// InfluxDBClient quản lý kết nối thực tế tới InfluxDB
type InfluxDBClient struct {
	Client influxdb2.Client
	Cfg    *config.Config
	Logger *logger.Logger
}

// NewInfluxDBClient tạo client mới
func NewInfluxDBClient(cfg *config.Config, logger *logger.Logger) *InfluxDBClient {
	client := influxdb2.NewClient(cfg.InfluxDBURL, cfg.InfluxDBAdminToken)
	return &InfluxDBClient{
		Client: client,
		Cfg:    cfg,
		Logger: logger,
	}
}

// WritePoint ghi 1 point vào InfluxDB
func (c *InfluxDBClient) WritePoint(measurement string, fields map[string]interface{}, tags map[string]string) error {
	writeAPI := c.Client.WriteAPIBlocking(c.Cfg.InfluxDBOrg, c.Cfg.InfluxDBBucket)
	point := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		c.Logger.Error("Error writing point: " + err.Error())
		return err
	}
	return nil
}

// Close đóng client
func (c *InfluxDBClient) Close() {
	c.Client.Close()
	c.Logger.Info("InfluxDB client closed")
}
