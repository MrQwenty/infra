package clients

import (
	"os"
	"github.com/coneno/logger"
	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
	emailAPI "github.com/influenzanet/messaging-service/pkg/api/email_client_service"
	whatsappAPI "github.com/influenzanet/messaging-service/pkg/api/whatsapp_client_service"
	"github.com/influenzanet/messaging-service/pkg/types"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc"
)

func connectToGRPCServer(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error.Fatalf("failed to connect to %s: %v", addr, err)
	}
	return conn
}

func ConnectToUserManagementService(addr string) (client umAPI.UserManagementApiClient, close func() error) {
	serverConn := connectToGRPCServer(addr)
	return umAPI.NewUserManagementApiClient(serverConn), serverConn.Close
}

func ConnectToEmailClientService(addr string) (client emailAPI.EmailClientServiceApiClient, close func() error) {
	serverConn := connectToGRPCServer(addr)
	return emailAPI.NewEmailClientServiceApiClient(serverConn), serverConn.Close
}

func ConnectToStudyService(addr string) (client studyAPI.StudyServiceApiClient, close func() error) {
	serverConn := connectToGRPCServer(addr)
	return studyAPI.NewStudyServiceApiClient(serverConn), serverConn.Close
}

func ConnectToLoggingService(addr string) (client loggingAPI.LoggingServiceApiClient, close func() error) {
	// Connect to user management service
	serverConn := connectToGRPCServer(addr)
	return loggingAPI.NewLoggingServiceApiClient(serverConn), serverConn.Close
}

func ConnectToWhatsAppClientService(addr string) (client whatsappAPI.WhatsAppClientServiceApiClient, close func() error) {
	serverConn := connectToGRPCServer(addr)
	return whatsappAPI.NewWhatsAppClientServiceApiClient(serverConn), serverConn.Close
}

func NewAPIClients() *types.APIClients {
	userMgmtURL := os.Getenv("USER_MANAGEMENT_SERVICE_URL")
	if userMgmtURL == "" {
		userMgmtURL = "user-management-service:5002"
	}
	
	emailServiceURL := os.Getenv("EMAIL_CLIENT_SERVICE_URL")
	if emailServiceURL == "" {
		emailServiceURL = "email-client-service:5005"
	}
	
	whatsappServiceURL := os.Getenv("WHATSAPP_CLIENT_SERVICE_URL")
	if whatsappServiceURL == "" {
		whatsappServiceURL = "whatsapp-client-service:5007"
	}
	
	studyServiceURL := os.Getenv("STUDY_SERVICE_URL")
	if studyServiceURL == "" {
		studyServiceURL = "study-service:5003"
	}
	
	loggingServiceURL := os.Getenv("LOGGING_SERVICE_URL")
	if loggingServiceURL == "" {
		loggingServiceURL = "logging-service:5006"
	}
	
	userMgmtClient, _ := ConnectToUserManagementService(userMgmtURL)
	emailClient, _ := ConnectToEmailClientService(emailServiceURL)
	whatsappClient, _ := ConnectToWhatsAppClientService(whatsappServiceURL)
	studyClient, _ := ConnectToStudyService(studyServiceURL)
	loggingClient, _ := ConnectToLoggingService(loggingServiceURL)
	
	return &types.APIClients{
		UserManagementService: userMgmtClient,
		EmailClientService:    emailClient,
		WhatsAppClientService: whatsappClient,
		StudyService:          studyClient,
		LoggingService:        loggingClient,
	}
}
