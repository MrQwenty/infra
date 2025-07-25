# Introduction 
This is the InfluenzaNet platform infrastructure with WhatsApp phone number verification system.

## WhatsApp Phone Number Verification

The system now includes a comprehensive phone number verification feature using WhatsApp Meta Business API:

### Features
- **Multi-method verification**: Support for both WhatsApp and SMS verification
- **Retry mechanism**: Automatic retry for failed message deliveries
- **Attempt limiting**: Users have a limited number of verification attempts (3 by default)
- **Expiration handling**: Verification codes expire after 10 minutes
- **Real-time feedback**: Live countdown timer and attempt counter
- **Graceful fallback**: Automatic cleanup of failed verifications
- **Meta Business API**: Direct integration with WhatsApp Graph API v19.0

### User Flow
1. User adds or changes phone number
2. System presents verification method selection (WhatsApp/SMS)
3. System sends verification code via selected method
4. User enters 6-digit code within 10 minutes
5. System verifies code with maximum 3 attempts
6. Phone number is marked as verified upon success

### WhatsApp API Configuration
The system is configured with the following Meta Business API credentials:
- **API URL**: https://graph.facebook.com/v19.0
- **Phone Number ID**: 676124925591256
- **Access Token**: Configured in environment variables
- **Template**: Uses 'hello_world' template for verification messages

### Technical Implementation
- **Frontend**: React components with TypeScript
- **Backend**: Go microservices with gRPC
- **WhatsApp Client**: Direct integration with Meta Business API
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
WHATSAPP_API_TOKEN=EAA6YlMvSeosBPM6NEv0SDKPu5IrTuAkLqL3jsQPglG181ZBD2bLy9P0TEtFJZBr064A5PFSc3fZAuZCMJeUKCkYbs2CN0vJkkuRLPJXqbaP8bpuX3ZC3PtX1yh7ZCexyjgTSjTy6PVUujRQ9cydJ4XV8ZBxYGRojZArolTo14YBnCgoJaf2VUdZBy1zOsQ4AyF1JrKD3gB64w8pBvhSTN1sDPyVjsSsoOiM25FyerpVdsqBVHaAZDZD
WHATSAPP_PHONE_NUMBER_ID=676124925591256
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

### Quick Test
Run the test script to verify WhatsApp API integration:
```bash
./script/test-whatsapp.sh
```

### Manual Testing
1. **Add Phone Number**: Use the settings page to add a phone number
2. **Select Method**: Choose WhatsApp or SMS verification
3. **Enter Code**: Input the 6-digit verification code
4. **Verify**: System validates the code and marks phone as verified

### API Testing
Test the WhatsApp API directly:
```bash
curl -X POST "https://graph.facebook.com/v19.0/676124925591256/messages" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "+1234567890",
    "type": "template",
    "template": {
      "name": "hello_world",
      "language": {"code": "en_US"}
    }
  }'
```

## API Endpoints

### Phone Number Management

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
