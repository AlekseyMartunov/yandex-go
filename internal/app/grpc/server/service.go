package server

import (
	"context"
	"fmt"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/grpc"
)

type encoder interface {
	Encode(url, userID string) (string, error)
	Decode(string) (string, error)
	GetAllURL(userID string) ([][2]string, error)
}

type ServiceURL struct {
	grpc.UnimplementedUrlServiceServer
	storage encoder
}

func NewServiceURL(e encoder) *ServiceURL {
	return &ServiceURL{
		storage: e,
	}
}

func (s *ServiceURL) Add(ctx context.Context, req *grpc.AddUrlRequest) (*grpc.AddUrlResponse, error) {
	var response grpc.AddUrlResponse

	shorted, err := s.storage.Encode(req.Key, req.UserId)
	if err != nil {
		response.Error = fmt.Sprintf("Ошибка при довалении URL: %s", req.Key)
	}
	response.Shorted = shorted
	return &response, nil

}

func (s *ServiceURL) GetShorted(ctx context.Context, req *grpc.GetUrlRequest) (*grpc.GetUrlResponse, error) {
	var response grpc.GetUrlResponse

	val, err := s.storage.Decode(req.Original)
	if err != nil {
		val = ""
	}

	response.Shorted = val
	return &response, nil
}

func (s *ServiceURL) GetAllForUser(ctx context.Context, req *grpc.GetAllUrlRequest) (*grpc.GetAllUrlResponse, error) {
	var response grpc.GetAllUrlResponse

	strs, err := s.storage.GetAllURL(req.UserId)
	if err != nil {
		response.PairUrl = nil
		return &response, nil
	}

	pairs := make([]*grpc.PairUrl, len(strs))

	for i := range strs {
		pairs[i].Original = strs[i][0]
		pairs[i].Shorted = strs[i][1]

	}

	response.PairUrl = pairs
	return &response, nil
}
