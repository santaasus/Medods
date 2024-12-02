package adapter

import (
	service "Medods/auth_service/inner_layer/service/auth"
	controller "Medods/auth_service/rest/controller"
)

func (a *BaseAdapter) AuthAdapter() *controller.Controller {
	service := service.Service{UserRepository: a.Repository}
	controller := controller.Controller{Service: &service}
	return &controller
}
