package forms

import (
	"websiteMonitor/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TransformGetSiteRequestToSite(in *pb.IdRequest) (*string, error) {
	if in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "The field 'Id' is missing")
	}

	id := string(in.Id)

	return &id, nil
}
