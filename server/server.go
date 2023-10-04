package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"websiteMonitor/forms"
	"websiteMonitor/handler"
	"websiteMonitor/models"
	"websiteMonitor/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	handler *handler.Handlers
	pb.UnsafeWebsiteMonitorServiceServer
}

func NewSiteServer(handler *handler.Handlers) pb.WebsiteMonitorServiceServer {
	return &Server{handler: handler}
}

func (s *Server) CreateSite(c context.Context, in *pb.SiteRequest) (*pb.SiteResponse, error) {
	site, customErr := forms.TransformCreateSiteRequestToSite(in)

	if customErr != nil {
		return nil, customErr
	}

	sitePtr, err := s.handler.Site.CreateSite(site)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao criar site: %v", err)
	}

	return &pb.SiteResponse{
		Id:         int32(sitePtr.ID),
		Name:       sitePtr.Name,
		DateCreate: timestamppb.New(sitePtr.CreatedAt),
		DateUpdate: timestamppb.New(sitePtr.UpdatedAt),
		DateDelete: timestamppb.New(sitePtr.DeletedAt.Time),
	}, nil
}

func (s *Server) DeleteSite(c context.Context, in *pb.IdRequest) (*pb.ResponseMessage, error) {
	id, err := forms.TransformDeleteSiteRequestToSite(in)

	if err != nil {
		return nil, err
	}

	err = s.handler.Site.DeleteSite(*id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao deletar site: %v", err)
	}

	return &pb.ResponseMessage{
		Message: "site deletado: " + fmt.Sprintf("%d", in.Id),
	}, nil
}

func (s *Server) UpdateSite(ctx context.Context, in *pb.SiteEditRequest) (*pb.SiteResponse, error) {
	site, customErr := forms.TransformUpdateSiteRequestToSite(in)

	if customErr != nil {
		return nil, customErr
	}

	sitePtr, err := s.handler.Site.UpdateSite(site)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao Deletar site: %v", err)
	}

	return &pb.SiteResponse{
		Id:         int32(sitePtr.ID),
		Name:       sitePtr.Name,
		DateCreate: timestamppb.New(sitePtr.CreatedAt),
		DateUpdate: timestamppb.New(sitePtr.UpdatedAt),
		DateDelete: timestamppb.New(sitePtr.DeletedAt.Time),
	}, nil
}

func (s *Server) GetSite(ctx context.Context, in *pb.IdRequest) (*pb.SiteResponse, error) {
	id, err := forms.TransformGetSiteRequestToSite(in)

	if err != nil {
		return nil, err
	}

	sitePtr, err := s.handler.Site.GetSite(*id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao buscar sites: %v", err)
	}

	return &pb.SiteResponse{
		Id:         int32(sitePtr.ID),
		Name:       sitePtr.Name,
		DateCreate: timestamppb.New(sitePtr.CreatedAt),
		DateUpdate: timestamppb.New(sitePtr.UpdatedAt),
		DateDelete: timestamppb.New(sitePtr.DeletedAt.Time),
	}, nil
}

func (s *Server) GetAllSites(ctx context.Context, in *pb.EmptyRequest) (*pb.SitesResponse, error) {
	var sites []*models.Site

	sites, err := s.handler.Site.GetAllSites()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao buscar sites: %v", err)
	}

	var siteResponses []*pb.SiteResponse
	for _, site := range sites {
		siteResponse := &pb.SiteResponse{
			Id:         int32(site.ID),
			Name:       site.Name,
			DateCreate: timestamppb.New(site.CreatedAt),
			DateUpdate: timestamppb.New(site.UpdatedAt),
			DateDelete: timestamppb.New(site.DeletedAt.Time),
		}
		siteResponses = append(siteResponses, siteResponse)
	}

	return &pb.SitesResponse{
		Sites: siteResponses,
	}, nil
}

func (s *Server) AutoMigrate(ctx context.Context, in *pb.EmptyRequest) (*pb.ResponseMessage, error) {

	fmt.Println("AutoMigrate main")
	fmt.Println(s.handler)
	fmt.Println(s.handler.Site)

	if err := s.handler.Site.AutoMigrate(); err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao criar modelo")
	}

	return &pb.ResponseMessage{Message: "tabela criada com sucesso"}, nil
}

var stopChan chan struct{}
var isRunning bool

func minhaRotina(s *Server) {
	for {
		select {
		case <-stopChan:
			fmt.Println("Rotina parada.")
			isRunning = false
			return
		default:
			s.iniciarMonitoramento()
			fmt.Println("Rotina rodando...")
			time.Sleep(2 * time.Second)
		}
	}
}

func (s *Server) iniciarMonitoramento() {

	fmt.Println("Monitorando...")

	sites, err := s.handler.Site.GetAllSites()
	if err != nil {
		return
	}

	for _, site := range sites {
		testaSite(site.Name)
	}

	fmt.Println("")

}

func testaSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
	} else {
		fmt.Println("Site:", site, "está com problemas. Status code:", resp.StatusCode)
	}
}

func (s *Server) IniciarRotina(context.Context, *pb.EmptyRequest) (*pb.ResponseMessage, error) {

	if isRunning {
		return nil, status.Errorf(200, "A rotina já esta rodando")
	}

	stopChan = make(chan struct{})
	go minhaRotina(s)
	isRunning = true

	return &pb.ResponseMessage{Message: "Rotina iniciada com sucesso!"}, nil
}

// PararRotina implements pb.WebsiteMonitorServiceServer.
func (*Server) PararRotina(context.Context, *pb.EmptyRequest) (*pb.ResponseMessage, error) {
	if !isRunning {
		return nil, status.Errorf(200, "A rotina não esta rodando")
	}

	close(stopChan)

	return &pb.ResponseMessage{Message: "Rotina parada com sucesso!"}, nil
}
