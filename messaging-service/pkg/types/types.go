package types

import (
	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	userAPI "github.com/influenzanet/user-management-service/pkg/api"
)

type APIClients struct {
	EmailClientService    messageAPI.MessagingServiceApiClient
	UserManagementService userAPI.UserManagementApiClient
	LoggingService        loggingAPI.LoggingServiceApiClient
}

type EmailTemplate struct {
	MessageType     string
	StudyKey        string
	DefaultLanguage string
	HeaderOverrides HeaderOverrides
	Translations    []LocalizedTemplate
}

type LocalizedTemplate struct {
	Lang        string
	Subject     string
	TemplateDef string
}

type HeaderOverrides struct {
	From      string
	Sender    string
	ReplyTo   []string
	NoReplyTo bool
}

func (h HeaderOverrides) ToEmailClientAPI() *messageAPI.HeaderOverrides {
	return &messageAPI.HeaderOverrides{
		From:      h.From,
		Sender:    h.Sender,
		ReplyTo:   h.ReplyTo,
		NoReplyTo: h.NoReplyTo,
	}
}

type OutgoingEmail struct {
	MessageType     string
	To              []string
	HeaderOverrides HeaderOverrides
	Subject         string
	Content         string
	HighPrio        bool
}

func EmailTemplateFromAPI(template *messageAPI.EmailTemplate) EmailTemplate {
	translations := make([]LocalizedTemplate, len(template.Translations))
	for i, t := range template.Translations {
		translations[i] = LocalizedTemplate{
			Lang:        t.Lang,
			Subject:     t.Subject,
			TemplateDef: t.TemplateDef,
		}
	}

	return EmailTemplate{
		MessageType:     template.MessageType,
		StudyKey:        template.StudyKey,
		DefaultLanguage: template.DefaultLanguage,
		HeaderOverrides: HeaderOverrides{
			From:      template.HeaderOverrides.From,
			Sender:    template.HeaderOverrides.Sender,
			ReplyTo:   template.HeaderOverrides.ReplyTo,
			NoReplyTo: template.HeaderOverrides.NoReplyTo,
		},
		Translations: translations,
	}
}
