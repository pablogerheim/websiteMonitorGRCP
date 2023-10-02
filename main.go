package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"websiteMonitor/database"
	"websiteMonitor/models"
	"websiteMonitor/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	banco string
	pb.UnsafeWebsiteMonitorServiceServer
}

// CriaNovoSite implements pb.WebsiteMonitorServiceServer.
func (s *server) CriaNovoSite(c context.Context, in *pb.SiteRequest) (*pb.SiteResponse, error) {
	var site models.Site
	site.Name = in.Name

	const layout = "2006-01-02 15:04:05"
	str := "2023-09-28 20:15:05"

	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if database.DB == nil {
		return &pb.SiteResponse{
			Id:         1,
			Name:       "hello",
			DateCreate: timestamppb.New(t),
			DateUpdate: timestamppb.New(t),
			DateDelete: timestamppb.New(t),
		}, nil
	}

	if database.DB == nil {
		return nil, status.Errorf(codes.Internal, "Conexão com o banco de dados é nil")
	}

	fmt.Println("site", in.Name)
	fmt.Println("erro DB server", database.DB)
	database.DB = database.DB.Debug()

	if err := database.DB.Create(&site).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao criar site: %v", err)
	}

	return &pb.SiteResponse{
		Id:         int32(site.ID),
		Name:       site.Name,
		DateCreate: timestamppb.New(site.CreatedAt),
		DateUpdate: timestamppb.New(site.UpdatedAt),
		DateDelete: timestamppb.New(site.DeletedAt.Time),
	}, nil
}

// DeletaSite implements pb.WebsiteMonitorServiceServer.
func (s *server) DeletaSite(c context.Context, in *pb.IdRequest) (*pb.ResponseMessage, error) {
	var site models.Site

	if err := database.DB.Delete(&site, in.Id).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao deletar site: %v", err)
	}

	return &pb.ResponseMessage{
		Message: "site deletado: " + fmt.Sprintf("%d", in.Id),
	}, nil
}

// EditaSite implements pb.WebsiteMonitorServiceServer.
func (s *server) EditaSite(ctx context.Context, in *pb.SiteEditRequest) (*pb.SiteResponse, error) {
	var site models.Site

	idInt := int(in.Id)

	result := database.DB.First(&site, idInt)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "Site not found with ID: %v", idInt)
	}

	site.Name = in.Name

	result = database.DB.Model(&site).Updates(site)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update site: %v", result.Error)
	}

	return &pb.SiteResponse{
		Id:         int32(site.ID),
		Name:       site.Name,
		DateCreate: timestamppb.New(site.CreatedAt),
		DateUpdate: timestamppb.New(site.UpdatedAt),
		DateDelete: timestamppb.New(site.DeletedAt.Time),
	}, nil
}

// ExibeTodosSites implements pb.WebsiteMonitorServiceServer.
func (s *server) ExibeTodosSites(ctx context.Context, in *pb.EmptyRequest) (*pb.SitesResponse, error) {
	var sites []models.Site

	// Consulta para buscar todos os registros
	result := database.DB.Find(&sites)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao buscar sites: %v", result.Error)
	}

	// Convertendo registros em SitesResponse
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

func (s *server) AutoMigrate(ctx context.Context, in *pb.EmptyRequest) (*pb.ResponseMessage, error) {

	if err := database.DB.AutoMigrate(&models.Site{}).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Erro ao deletar site")
	}

	return &pb.ResponseMessage{Message: "tabela criada com sucesso"}, nil
}

type Site struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Status string `json:"status"`
}

var stopChan chan struct{}
var isRunning bool

func minhaRotina() {
	for {
		select {
		case <-stopChan:
			fmt.Println("Rotina parada.")
			isRunning = false
			return
		default:
			iniciarMonitoramento()
			fmt.Println("Rotina rodando...")
			time.Sleep(2 * time.Second)
		}
	}
}

func iniciarMonitoramento() {

	fmt.Println("Monitorando...")

	var sites []models.Site

	result := database.DB.Find(&sites)
	if result.Error != nil {
		return
	}

	fmt.Println("sites", sites)

	for _, site := range sites {
		fmt.Println("Testando Name:", site.Name)
		fmt.Println("Testando site:", site)
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

// IniciarRotina implements pb.WebsiteMonitorServiceServer.
func (*server) IniciarRotina(context.Context, *pb.EmptyRequest) (*pb.ResponseMessage, error) {

	if isRunning {
		return nil, status.Errorf(200, "A rotina já esta rodando")
	}

	stopChan = make(chan struct{})
	go minhaRotina()
	isRunning = true

	return &pb.ResponseMessage{Message: "Rotina iniciada com sucesso!"}, nil
}

// PararRotina implements pb.WebsiteMonitorServiceServer.
func (*server) PararRotina(context.Context, *pb.EmptyRequest) (*pb.ResponseMessage, error) {
	if !isRunning {
		return nil, status.Errorf(200, "A rotina não esta rodando")
	}
	
	close(stopChan)
	
	return &pb.ResponseMessage{Message: "Rotina parada com sucesso!"}, nil
}

func main() {
	fmt.Println("start server")
	database.ConectaComBancoDeDados()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWebsiteMonitorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
