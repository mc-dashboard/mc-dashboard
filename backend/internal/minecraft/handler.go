package minecraft

import (
	"net/http"

	"github.com/rohanvsuri/minecraft-dashboard/internal/lambda"
)

type MinecraftHandler struct {
	LambdaService *lambda.FunctionWrapper
}

func NewMinecraftHandler(lambdaService *lambda.FunctionWrapper) *MinecraftHandler {
	return &MinecraftHandler{LambdaService: lambdaService}
}

func (h *MinecraftHandler) StartServer(w http.ResponseWriter, r *http.Request) {
	h.LambdaService.CallLambda("ec2-start")
}

func (h *MinecraftHandler) StopServer(w http.ResponseWriter, r *http.Request) {
	h.LambdaService.CallLambda("ec2-stop")
}

// future implementations:
// func (h *Handler) RestartServer(w http.ResponseWriter, r *http.Request) {
// 	h.LambdaService.CallLambda("restartServer")
// }

// func (h *Handler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
// 	h.LambdaService.CallLambda("getServerStatus")
// }