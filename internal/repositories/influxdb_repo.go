package repository

import (
	influxdb_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/influxdb"
)

type MeasurementRepository interface {
	WriteMeasurement(measurement string, fields map[string]interface{}, tags map[string]string) error
}

type influxDBRepo struct {
	client *influxdb_client.InfluxDBClient
}

func NewInfluxDBRepo(client *influxdb_client.InfluxDBClient) MeasurementRepository {
	return &influxDBRepo{
		client: client,
	}
}

func (r *influxDBRepo) WriteMeasurement(measurement string, fields map[string]interface{}, tags map[string]string) error {
	return r.client.WritePoint(measurement, fields, tags)
}
