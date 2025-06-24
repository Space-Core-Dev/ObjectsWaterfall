package services

import "object-shooter.com/core/models"

type Send interface {
	SendRequest(host string, obj interface{}, headers map[string]string) (models.ResponseResult, error)
}
