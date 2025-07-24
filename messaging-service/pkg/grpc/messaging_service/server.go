package messaging_service

import (
	"net"
	"os"

	"github.com/coneno/logger"
	api "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	"github.com/influenzanet/messaging-service/pkg/dbs/messagedb"
	"github.com/influenzanet/messaging-service/pkg/types"
	"google.golang.org/grpc"
)

const apiVersion = "v1.0.0"

type messagingServer struct {
	api.UnimplementedMessagingServiceApiServer
	messageDBservice *messagedb.MessageDBService
	clients          *types.APIClients
}

func NewMessagingServer(
	messageDBservice *messagedb.MessageDBService,
	clients *types.APIClients,
) *messagingServer {
	return &messagingServer{
		messageDBservice: messageDBservice,
		clients:          clients,
	}
}

func RunServer() error {
	port := os.Getenv("MESSAGING_SERVICE_LISTEN_PORT")
	if port == "" {
		port = "5005"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	
	messageDBCon := os.Getenv("MESSAGE_DB_CONNECTION_STR")
	messageDBUsername := os.Getenv("MESSAGE_DB_USERNAME")
	messageDBPassword := os.Getenv("MESSAGE_DB_PASSWORD")
	messageDBConnectionPrefix := os.Getenv("MESSAGE_DB_CONNECTION_PREFIX")

	messageDB := messagedb.NewMessageDBService(
		messageDBCon,
		messageDBUsername,
		messageDBPassword,
		messageDBConnectionPrefix,
		"messagedb",
	)

	clients := &types.APIClients{}

	api.RegisterMessagingServiceApiServer(grpcServer, NewMessagingServer(messageDB, clients))

	logger.Info.Printf("Messaging service listening on port %s", port)
	return grpcServer.Serve(lis)
}
