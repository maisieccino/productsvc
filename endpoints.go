package productsvc

import (
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetProductEndpoint  endpoint.Endpoint
	PostProductEndpoint endpoint.Endpoint
	ListProductEndpoint endpoint.Endpoint
}
