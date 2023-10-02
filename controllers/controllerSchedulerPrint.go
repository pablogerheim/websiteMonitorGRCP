package controllers


import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

func IniciarRotina(c *gin.Context) {
	if isRunning {
		c.JSON(400, gin.H{"message": "A rotina já está rodando."})
		return
	}

	stopChan = make(chan struct{})
	go minhaRotina()
	isRunning = true
	c.JSON(200, gin.H{"message": "Rotina iniciada."})
}

func PararRotina(c *gin.Context) {
	if !isRunning {
		c.JSON(400, gin.H{"message": "A rotina não está rodando."})
		return
	}

	close(stopChan)
	c.JSON(200, gin.H{"message": "Rotina parada."})
}

func iniciarMonitoramento() {

	fmt.Println("Monitorando...")
	resp, err := http.Get("http://localhost:8080/site")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		return
	}

	fmt.Println("Resp", resp)
	fmt.Println("Resp Body", resp.Body)
	defer resp.Body.Close()
	var sites []Site
	if err := json.NewDecoder(resp.Body).Decode(&sites); err != nil {
		fmt.Println("Ocorreu um erro ao decodificar o body:", err)
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
