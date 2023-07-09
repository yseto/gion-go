package handler

type ApiServer struct{}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

var _ StrictServerInterface = (*ApiServer)(nil)
