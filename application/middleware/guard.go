package guard

import (
	"antrein/dd-dashboard-config/model/config"
	"antrein/dd-dashboard-config/model/dto"
	"antrein/dd-dashboard-config/model/entity"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/antrein/proto-repository/pb/dd"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GuardContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type AuthGuardContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Claims         entity.JWTClaim
}

func (g *GuardContext) ReturnError(status int, message string) error {
	g.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  status,
		Message: message,
	})
}

func (g *GuardContext) ReturnSuccess(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func (g *GuardContext) ReturnCreated(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Created",
		Data:    data,
	})
}

func (g *GuardContext) ReturnEvent(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(g.ResponseWriter, "data: %s\n\n", jsonData)
	if err != nil {
		return err // Handle writing errors
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		return fmt.Errorf("streaming unsupported")
	}

	return nil
}

func (g *AuthGuardContext) ReturnError(status int, message string) error {
	g.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  status,
		Message: message,
	})
}

func (g *AuthGuardContext) ReturnSuccess(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func (g *AuthGuardContext) ReturnCreated(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Created",
		Data:    data,
	})
}

func (g *AuthGuardContext) ReturnEvent(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(g.ResponseWriter, "data: %s\n\n", jsonData)
	if err != nil {
		return err // Handle writing errors
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		return fmt.Errorf("streaming unsupported")
	}

	return nil
}

func DefaultGuard(handlerFunc func(g *GuardContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guardCtx := GuardContext{
			ResponseWriter: w,
			Request:        r,
		}
		if err := handlerFunc(&guardCtx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func AuthGuard(cfg *config.Config, handlerFunc func(g *AuthGuardContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		grpcClient, err := grpc.Dial(cfg.GRPCConfig.DashboardAuth, grpc.WithTransportCredentials(insecure.NewCredentials()))

		client := pb.NewAuthServiceClient(grpcClient)
		ctx := context.Background()
		authResp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: authHeader})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !authResp.IsValid {
			http.Error(w, "Token invalid", http.StatusUnauthorized)
			return
		}

		authGuardCtx := AuthGuardContext{
			ResponseWriter: w,
			Request:        r,
			Claims: entity.JWTClaim{
				UserID: authResp.UserId,
			},
		}

		if err := handlerFunc(&authGuardCtx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func BodyParser(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func IsMethod(r *http.Request, method string) bool {
	return r.Method == method
}

func GetParam(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}
