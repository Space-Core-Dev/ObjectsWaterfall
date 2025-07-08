package mappers

import (
	"objectswaterfall.com/application/dtos"
	"objectswaterfall.com/core/models"
)

func FromDtoToWorkerSettings(dto dtos.BackgroundWorkerSettingsDto) models.BackgroundWorkerSettings {
	return models.BackgroundWorkerSettings{
		TableName:          dto.TableName,
		Timer:              dto.Timer,
		RequestDelay:       dto.RequestDelay,
		Random:             dto.Random,
		WritesNumberToSend: dto.WritesNumberToSend,
		TotalToSend:        dto.TotalToSend,
		StopWhenTableEnds:  dto.StopWhenTableEnds,
	}
}
