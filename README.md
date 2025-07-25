# Introduction 
This is the InfluenzaNet platform infrastructure with WhatsApp phone number verification system.

## WhatsApp Phone Number Verification

The system now includes a comprehensive phone number verification feature using WhatsApp:

### Features
- **Multi-method verification**: Support for both WhatsApp and SMS verification
- **Retry mechanism**: Automatic retry for failed message deliveries
- **Attempt limiting**: Users have a limited number of verification attempts (3 by default)
- **Expiration handling**: Verification codes expire after 10 minutes
- **Real-time feedback**: Live countdown timer and attempt counter
- **Graceful fallback**: Automatic cleanup of failed verifications

### User Flow
1. User adds or changes phone number
2. System presents verification method selection (WhatsApp/SMS)
3. System sends verification code via selected method
4. User enters 6-digit code within 10 minutes
5. System verifies code with maximum 3 attempts
6. Phone number is marked as verified upon success

### Technical Implementation
- **Frontend**: React components with TypeScript
- **Backend**: Go microservices with gRPC
- **Verification Service**: Dedicated WhatsApp verification service
- **Retry Logic**: Exponential backoff for failed deliveries
- **Security**: Token-based verification with expiration
- **Cleanup**: Automatic cleanup of expired verification attempts

# Getting Started

## Prerequisites
- Docker and Docker Compose
- Node.js and Yarn (for frontend development)
- Go 1.19+ (for backend development)

## Installation

1. Clone the repository and run the initialization script:
```bash
./script/init-repos.sh
```

2. Configure WhatsApp API credentials in `.env`:
```bash
WHATSAPP_API_TOKEN=your_whatsapp_api_token_here
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id_here
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_webhook_verify_token_here
```

3. Start the services:
```bash
docker-compose up -d
```

4. Import sample data (optional):
```bash
./script/import-db.sh
```

## Development

### Frontend Development
```bash
cd participant-webapp
yarn install
yarn start
```

### Backend Development
To rebuild specific services:
```bash
./script/rebuild-ms-influenza.sh
```

# Build and Test

## Testing WhatsApp Verification

1. **Add Phone Number**: Use the settings page to add a phone number
2. **Select Method**: Choose WhatsApp or SMS verification
3. **Enter Code**: Input the 6-digit verification code
4. **Verify**: System validates the code and marks phone as verified

## API Endpoints

### Phone Number Management
- `POST /v1/user/phone/add` - Add phone number
- `POST /v1/user/phone/change` - Change phone number
- `POST /v1/user/whatsapp-verification/verify` - Verify code
- `POST /v1/user/whatsapp-verification/resend` - Resend code
- `POST /v1/user/whatsapp-verification/cancel` - Cancel verification

# Contribute

## Contributing Guidelines

1. **Code Style**: Follow existing patterns and conventions
2. **Testing**: Add tests for new verification features
3. **Documentation**: Update documentation for API changes
4. **Security**: Ensure proper validation and error handling

## Architecture

### Verification Flow
```
Frontend → API Gateway → User Management Service → WhatsApp Client Service
```

### Key Components
- **VerifyWhatsApp.tsx**: Main verification dialog component
- **whatsappVerificationService.ts**: Frontend verification service
- **account_management_endpoints.go**: Backend verification handlers
- **whatsapp-client-service**: Dedicated WhatsApp messaging service

### Security Features
- Token-based verification with expiration
- Rate limiting for verification attempts
- Phone number format validation
- Automatic cleanup of expired sessions
- Secure code generation and validation
