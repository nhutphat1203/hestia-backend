package jobs

import (
	"github.com/nhutphat1203/hestia-backend/internal/domain"
	service "github.com/nhutphat1203/hestia-backend/internal/services"
)

type SaveMeasurementJob struct {
	SensorData *domain.SensorData
	Service    *service.MeasurementService
}

func NewSaveMeasurementJob(sensorData domain.SensorData, service *service.MeasurementService) *SaveMeasurementJob {
	return &SaveMeasurementJob{
		SensorData: &sensorData,
		Service:    service,
	}
}

func (j *SaveMeasurementJob) Execute() error {
	return j.Service.RecordData(j.SensorData)
}
