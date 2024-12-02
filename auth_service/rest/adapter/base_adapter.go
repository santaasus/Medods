package adapter

import repository "Medods/auth_service/inner_layer/repository/user"

type BaseAdapter struct {
	Repository repository.IRepository
}
