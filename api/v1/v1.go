// v1版本的 HTTP 接口
package v1

import "zgoframe/service"

var apiServicesMount *service.MyService

func ApiServices() *service.MyService {
	if apiServicesMount == nil {
		apiServicesMount = service.NewMyService()
	}
	return apiServicesMount
}
