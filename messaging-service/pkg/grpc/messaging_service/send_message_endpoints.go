package messaging_service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/coneno/logger"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/influenzanet/go-utils/pkg/constants"
	"github.com/influenzanet/go-utils/pkg/token_checks"
	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
	emailAPI "github.com/influenzanet/messaging-service/pkg/api/email_client_service"
	api "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	"github.com/influenzanet/messaging-service/pkg/bulk_messages"
	"github.com/influenzanet/messaging-service/pkg/templates"
	"github.com/influenzanet/messaging-service/pkg/types"
	userAPI "github.com/influenzanet/user-management-service/pkg/api"
	"github.com/influenzanet/go-utils/pkg/api_types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *messagingServer) Status(ctx context.Context, _ *empty.Empty) (*api.ServiceStatus, error) {
	return &api.ServiceStatus{
		Status:  api.ServiceStatus_NORMAL,
		Msg:     "service running",
		Version: apiVersion,
	}, nil
}

func (s *messagingServer) SendMessageToAllUsers(ctx context.Context, req *api.SendMessageToAllUsersReq) (*api.ServiceStatus, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.Template == nil {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if !token_checks.CheckIfAnyRolesInToken(req.Token, []string{constants.USER_ROLE_ADMIN}) {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_BULK_MESSAGE_SEND, fmt.Sprintf("permission denied for send %s to all users", req.Template.MessageType))
		return nil, status.Error(codes.PermissionDenied, "no permission to send messages")
	}

	// use go method (don't wait for result since it can take long)
	go bulk_messages.GenerateForAllUsers(
		s.clients,
		s.messageDBservice,
		req.Token.InstanceId,
		types.EmailTemplateFromAPI(req.Template),
		req.IgnoreWeekday,
		"one time message",
	)
	return &api.ServiceStatus{
		Msg:     "message sending triggered",
		Status:  api.ServiceStatus_NORMAL,
		Version: apiVersion,
	}, nil
}

func (s *messagingServer) SendMessageToStudyParticipants(ctx context.Context, req *api.SendMessageToStudyParticipantsReq) (*api.ServiceStatus, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" || req.Template == nil {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}
	if !token_checks.CheckIfAnyRolesInToken(req.Token, []string{constants.USER_ROLE_RESEARCHER, constants.USER_ROLE_ADMIN}) {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_BULK_MESSAGE_SEND, fmt.Sprintf("permission denied for send %s to study %s", req.Template.MessageType, req.StudyKey))
		return nil, status.Error(codes.PermissionDenied, "no permission to send messages")
	}
	req.Template.StudyKey = req.StudyKey

	// use go method (don't wait for result since it can take long)
	go bulk_messages.GenerateForStudyParticipants(
		s.clients,
		s.messageDBservice,
		req.Token.InstanceId,
		types.EmailTemplateFromAPI(req.Template),
		req.Condition,
		req.IgnoreWeekday,
		"one time message",
	)
	return &api.ServiceStatus{
		Msg:     "message sending triggered",
		Status:  api.ServiceStatus_NORMAL,
		Version: apiVersion,
	}, nil
}

func (s *messagingServer) SendInstantEmail(ctx context.Context, req *api.SendEmailReq) (*api.ServiceStatus, error) {
	if req == nil || req.InstanceId == "" || len(req.To) < 1 || req.MessageType == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	templateDef, err := s.messageDBservice.FindEmailTemplateByType(req.InstanceId, req.MessageType, req.StudyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "template not found")
	}

	translation := templates.GetTemplateTranslation(templateDef, req.PreferredLanguage)

	decodedTemplate, err := base64.StdEncoding.DecodeString(translation.TemplateDef)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req.ContentInfos == nil {
		req.ContentInfos = map[string]string{}
	}
	globalTemplateInfos := templates.LoadGlobalEmailTemplateConstants()
	for k, v := range globalTemplateInfos {
		req.ContentInfos[k] = v
	}

	req.ContentInfos["language"] = req.PreferredLanguage
	// execute template
	templateName := req.InstanceId + req.MessageType + req.PreferredLanguage
	content, err := templates.ResolveTemplate(
		templateName,
		string(decodedTemplate),
		req.ContentInfos,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "content could not be generated")
	}

	outgoingEmail := types.OutgoingEmail{
		MessageType:     req.MessageType,
		To:              req.To,
		HeaderOverrides: templateDef.HeaderOverrides,
		Subject:         translation.Subject,
		Content:         content,
		HighPrio:        !req.UseLowPrio,
	}

	_, err = s.clients.EmailClientService.SendEmail(ctx, &emailAPI.SendEmailReq{
		To:              outgoingEmail.To,
		HeaderOverrides: outgoingEmail.HeaderOverrides.ToEmailClientAPI(),
		Subject:         outgoingEmail.Subject,
		Content:         content,
		HighPrio:        !req.UseLowPrio,
	})
	if err != nil {
		_, errS := s.messageDBservice.AddToOutgoingEmails(req.InstanceId, outgoingEmail)
		if errS != nil {
			logger.Error.Printf("Error while saving to outgoing: %v", errS)
		}
		return &api.ServiceStatus{
			Version: apiVersion,
			Msg:     "failed sending message, added to outgoing",
			Status:  api.ServiceStatus_PROBLEM,
		}, nil
	}

	_, err = s.messageDBservice.AddToSentEmails(req.InstanceId, outgoingEmail)
	if err != nil {
		logger.Error.Printf("Saving to sent: %v", err)
	}

	return &api.ServiceStatus{
		Version: apiVersion,
		Msg:     "message sent",
		Status:  api.ServiceStatus_NORMAL,
	}, nil
}

func (s *messagingServer) QueueEmailTemplateForSending(ctx context.Context, req *api.SendEmailReq) (*api.ServiceStatus, error) {
	if req == nil || req.InstanceId == "" || len(req.To) < 1 || req.MessageType == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	templateDef, err := s.messageDBservice.FindEmailTemplateByType(req.InstanceId, req.MessageType, req.StudyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "template not found")
	}

	translation := templates.GetTemplateTranslation(templateDef, req.PreferredLanguage)

	decodedTemplate, err := base64.StdEncoding.DecodeString(translation.TemplateDef)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req.ContentInfos == nil {
		req.ContentInfos = map[string]string{}
	}
	globalTemplateInfos := templates.LoadGlobalEmailTemplateConstants()
	for k, v := range globalTemplateInfos {
		req.ContentInfos[k] = v
	}

	req.ContentInfos["language"] = req.PreferredLanguage
	// execute template
	templateName := req.InstanceId + req.MessageType + req.PreferredLanguage
	content, err := templates.ResolveTemplate(
		templateName,
		string(decodedTemplate),
		req.ContentInfos,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "content could not be generated")
	}

	outgoingEmail := types.OutgoingEmail{
		MessageType:     req.MessageType,
		To:              req.To,
		HeaderOverrides: templateDef.HeaderOverrides,
		Subject:         translation.Subject,
		Content:         content,
		HighPrio:        !req.UseLowPrio,
	}

	_, err = s.messageDBservice.AddToOutgoingEmails(req.InstanceId, outgoingEmail)
	if err != nil {
		logger.Error.Printf("Error while saving to outgoing: %v", err)
		return &api.ServiceStatus{
			Version: apiVersion,
			Msg:     "failed adding message to outgoing",
			Status:  api.ServiceStatus_PROBLEM,
		}, nil
	}

	return &api.ServiceStatus{
		Version: apiVersion,
		Msg:     "message added to ougoing",
		Status:  api.ServiceStatus_NORMAL,
	}, nil
}

func (s *messagingServer) SendNotificationToPreferredChannels(ctx context.Context, req *api.SendEmailReq) (*api.ServiceStatus, error) {
	if req == nil || req.InstanceId == "" || len(req.To) < 1 || req.MessageType == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	var successCount int
	var errors []string

	for _, userID := range req.To {
		user, err := s.clients.UserManagementService.GetUser(ctx, &userAPI.UserReference{
			Token:  &api_types.TokenInfos{InstanceId: req.InstanceId},
			UserId: userID,
		})
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to get user %s: %v", userID, err))
			continue
		}

		if user == nil || user.ContactPreferences == nil {
			errors = append(errors, fmt.Sprintf("user %s has no contact preferences", userID))
			continue
		}

		prefs := user.ContactPreferences
		channels := prefs.PreferredChannels
		if len(channels) == 0 {
			channels = []string{"email"}
		}

		var channelSuccess bool
		for _, channel := range channels {
			switch channel {
			case "email":
				emailReq := &api.SendEmailReq{
					InstanceId:        req.InstanceId,
					To:                []string{user.Account.AccountId},
					MessageType:       req.MessageType,
					StudyKey:          req.StudyKey,
					PreferredLanguage: req.PreferredLanguage,
					ContentInfos:      req.ContentInfos,
					UseLowPrio:        req.UseLowPrio,
				}
				_, err := s.QueueEmailTemplateForSending(ctx, emailReq)
				if err == nil {
					channelSuccess = true
					logger.Debug.Printf("Email notification queued for user %s", userID)
				} else {
					logger.Error.Printf("Failed to queue email for user %s: %v", userID, err)
				}

			case "whatsapp":
				if !prefs.SubscribedToWhatsApp || prefs.WhatsAppNumber == "" {
					logger.Debug.Printf("User %s not subscribed to WhatsApp or no number", userID)
					continue
				}

				templateDef, err := s.messageDBservice.FindEmailTemplateByType(req.InstanceId, req.MessageType, req.StudyKey)
				if err != nil {
					logger.Error.Printf("Template not found for WhatsApp message: %v", err)
					continue
				}

				translation := templates.GetTemplateTranslation(templateDef, req.PreferredLanguage)
				decodedTemplate, err := base64.StdEncoding.DecodeString(translation.TemplateDef)
				if err != nil {
					logger.Error.Printf("Failed to decode template: %v", err)
					continue
				}

				contentInfos := req.ContentInfos
				if contentInfos == nil {
					contentInfos = map[string]string{}
				}
				globalTemplateInfos := templates.LoadGlobalEmailTemplateConstants()
				for k, v := range globalTemplateInfos {
					contentInfos[k] = v
				}
				contentInfos["language"] = req.PreferredLanguage

				templateName := req.InstanceId + req.MessageType + req.PreferredLanguage + "_whatsapp"
				content, err := templates.ResolveTemplate(templateName, string(decodedTemplate), contentInfos)
				if err != nil {
					logger.Error.Printf("Failed to resolve WhatsApp template: %v", err)
					continue
				}

				content = strings.ReplaceAll(content, "<br>", "\n")
				content = strings.ReplaceAll(content, "<br/>", "\n")
				content = strings.ReplaceAll(content, "<p>", "")
				content = strings.ReplaceAll(content, "</p>", "\n")

				whatsappContent := content
				maxRetries := 5
				highPrio := !req.UseLowPrio
				baseDelay := 30
				
				if req.MessageType == "whatsapp_verification" {
					code := req.ContentInfos["code"]
					whatsappContent = fmt.Sprintf("Il tuo codice di verifica InfluenzaNet è: %s", code)
					maxRetries = 3
					highPrio = true
					baseDelay = 15
				}
				
				outgoingWhatsApp := types.OutgoingWhatsApp{
					MessageType:      req.MessageType,
					To:              prefs.WhatsAppNumber,
					Content:         whatsappContent,
					HighPrio:        highPrio,
					MaxRetries:      maxRetries,
					BaseDelaySeconds: baseDelay,
					NextRetryAt:     0,
				}

				_, err = s.messageDBservice.AddToOutgoingWhatsApp(req.InstanceId, outgoingWhatsApp)
				if err == nil {
					channelSuccess = true
					logger.Debug.Printf("WhatsApp notification queued for user %s", userID)
				} else {
					logger.Error.Printf("Failed to queue WhatsApp for user %s: %v", userID, err)
				}
			}
		}

		if channelSuccess {
			successCount++
		} else {
			errors = append(errors, fmt.Sprintf("failed to send notification to user %s via any preferred channel", userID))
		}
	}

	if successCount == 0 {
		return &api.ServiceStatus{
			Version: apiVersion,
			Msg:     fmt.Sprintf("failed to send notifications: %s", strings.Join(errors, "; ")),
			Status:  api.ServiceStatus_PROBLEM,
		}, nil
	}

	msg := fmt.Sprintf("notifications sent to %d/%d users", successCount, len(req.To))
	if len(errors) > 0 {
		msg += fmt.Sprintf(" (errors: %s)", strings.Join(errors, "; "))
	}

	return &api.ServiceStatus{
		Version: apiVersion,
		Msg:     msg,
		Status:  api.ServiceStatus_NORMAL,
	}, nil
}
