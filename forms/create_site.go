package forms

import (
	"websiteMonitor/models"
	"websiteMonitor/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TransformCreateSiteRequestToSite(in *pb.SiteRequest) (*models.Site, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "The field 'name' is missing")
	}
	site := new(models.Site)
	site.Name = in.Name
	return site, nil
}
