# WhatsApp Verification Implementation Guide

## Overview
This document outlines the backend implementation requirements for WhatsApp verification in changephone and addphone functions.

## Infrastructure Changes Completed
- Added `whatsapp-client-service` to docker-compose.yml
- Updated messaging-service and message-scheduler to include WhatsApp client
- Added WhatsApp configuration files and environment variables
- Updated frontend dialogs to support verification method selection

## Backend Changes Required

### 1. Messaging Service (../messaging-service)
Create a new WhatsApp client service similar to the existing email-client-service:

```
messaging-service/
├── build/docker/whatsapp-client-service/
│   └── Dockerfile
├── cmd/whatsapp-client-service/
│   └── main.go
└── pkg/whatsapp/
    ├── client.go
    ├── templates.go
    └── verification.go
```

Key components:
- WhatsApp Business API integration using Facebook Graph API
- Message template management for verification codes
- Error handling and retry logic
- Logging and monitoring integration

### 2. User Management Service (../user-management-service)
Modify the contact management endpoints to support verification method parameter:

#### API Changes for /v1/user/contact/add-phone
```json
// Request
{
  "phone": "+393930238999",
  "verificationMethod": "whatsapp" // or "sms"
}

// Response
{
  "success": true,
  "verificationMethod": "whatsapp",
  "message": "WhatsApp verification sent to +393930238999",
  "verificationId": "uuid-here"
}
```

#### API Changes for /v1/user/contact/change-phone
```json
// Request
{
  "newPhone": "+393930238999",
  "verificationMethod": "whatsapp" // or "sms"
}

// Response
{
  "success": true,
  "verificationMethod": "whatsapp",
  "message": "WhatsApp verification sent to +393930238999",
  "verificationId": "uuid-here"
}
```

### 3. API Gateway (../api-gateway)
Update participant-api to handle WhatsApp verification requests:
- Add validation for verificationMethod parameter
- Route WhatsApp verification requests to appropriate service
- Handle WhatsApp-specific error responses

### 4. Frontend Integration (participant-webapp)
The frontend needs to be updated to:
- Add verification method selection UI in phone dialogs
- Send verificationMethod parameter in API requests
- Handle WhatsApp-specific success/error messages
- Update form validation to include verification method

## WhatsApp Business API Integration

### Required Setup
1. Facebook Business Account
2. WhatsApp Business API access
3. Phone number verification with Meta
4. Message templates approval

### Message Templates
Create and get approval for these templates:
- Phone verification code (English)
- Phone verification code (Italian)
- Verification failure notification
- Verification success confirmation

### Environment Variables
```bash
WHATSAPP_API_TOKEN=your_token_here
WHATSAPP_PHONE_NUMBER_ID=your_phone_id_here
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_webhook_token_here
```

## Implementation Steps

### Phase 1: WhatsApp Client Service
1. Create WhatsApp client service in messaging-service repository
2. Implement Facebook Graph API integration
3. Add message template management
4. Create Docker configuration
5. Add unit tests

### Phase 2: User Management Updates
1. Update contact endpoints to accept verificationMethod parameter
2. Integrate with WhatsApp client service
3. Update verification token handling
4. Add database schema changes if needed
5. Update API documentation

### Phase 3: Frontend Updates
1. Update phone dialogs to include verification method selection
2. Modify API calls to include verificationMethod parameter
3. Update error handling for WhatsApp-specific errors
4. Add UI tests for new verification flow

### Phase 4: Testing and Deployment
1. Integration testing with WhatsApp Business API
2. End-to-end testing of phone verification flows
3. Load testing for verification message sending
4. Production deployment with monitoring

## Error Handling

### WhatsApp-Specific Errors
- Invalid phone number format
- WhatsApp not available for phone number
- Message template not approved
- API rate limiting
- Service unavailable

### Fallback Strategy
If WhatsApp verification fails, the system should:
1. Log the error
2. Automatically fallback to SMS verification
3. Notify the user about the fallback
4. Continue with SMS verification flow

## Monitoring and Logging

### Metrics to Track
- WhatsApp verification success rate
- Message delivery time
- API error rates
- User preference for verification method
- Fallback to SMS frequency

### Logging Requirements
- All WhatsApp API calls
- Verification attempts and results
- Error conditions and fallbacks
- User verification method preferences

## Security Considerations

### Data Protection
- Encrypt WhatsApp API tokens
- Secure webhook endpoints
- Validate all incoming webhook data
- Rate limiting for verification requests

### Privacy Compliance
- User consent for WhatsApp communication
- Data retention policies for verification logs
- GDPR compliance for EU users
- Option to opt-out of WhatsApp verification

## Testing Strategy

### Unit Tests
- WhatsApp client service functions
- Message template rendering
- Error handling scenarios
- API response parsing

### Integration Tests
- End-to-end verification flow
- WhatsApp API integration
- Database operations
- Service communication

### Manual Testing
- Real phone number verification
- Different phone number formats
- International phone numbers
- Error scenarios and fallbacks

## Deployment Checklist

- [ ] WhatsApp Business API account setup
- [ ] Message templates approved by Meta
- [ ] Environment variables configured
- [ ] Docker images built and tested
- [ ] Database migrations applied
- [ ] Monitoring and alerting configured
- [ ] Documentation updated
- [ ] Team training completed

## Support and Maintenance

### Ongoing Tasks
- Monitor WhatsApp API changes
- Update message templates as needed
- Review verification success rates
- Optimize message delivery performance
- Handle user feedback and issues

### Escalation Procedures
- WhatsApp API outages
- Message delivery failures
- Security incidents
- Performance degradation
