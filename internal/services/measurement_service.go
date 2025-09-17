package service

import (
	"github.com/nhutphat1203/hestia-backend/internal/domain"
	repository "github.com/nhutphat1203/hestia-backend/internal/repositories"
)

type MeasurementService struct {
	repo repository.MeasurementRepository
}

func NewMeasurementService(repo repository.MeasurementRepository) *MeasurementService {
	return &MeasurementService{repo: repo}
}

func (s *MeasurementService) RecordData(data *domain.SensorData) error {
	fields := map[string]interface{}{
		"temperature": data.Measure.T,
		"humidity":    data.Measure.H,
		"pressure":    data.Measure.P,
		"lux":         data.Measure.Lux,
		"score":       data.Score,
		"state":       data.State,
	}

	tags := map[string]string{
		"room":   data.RoomID,
		"type":   data.Type,
		"source": data.Meta.Source,
	}
	return s.repo.WriteMeasurement("environment", fields, tags)
}
