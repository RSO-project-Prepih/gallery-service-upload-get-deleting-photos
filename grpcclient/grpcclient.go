package grpcclient

import (
	"log"

	get_photo_info "github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/github.com/RSO-project-Prepih/get-photo-info"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPhotoServiceClient(address string) (get_photo_info.PhotoServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to dial gRPC server: %v", err)
		return nil, err
	}

	return get_photo_info.NewPhotoServiceClient(conn), nil
}
