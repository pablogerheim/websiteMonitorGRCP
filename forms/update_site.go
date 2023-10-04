package forms

import (
	"websiteMonitor/models"
	"websiteMonitor/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TransformUpdateSiteRequestToSite(in *pb.SiteEditRequest) (*models.Site, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "The field 'name' is missing")
	}
		if in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "The field 'Id' is missing")
	}
	site := new(models.Site)
	site.Name = in.Name
	return site, nil
}