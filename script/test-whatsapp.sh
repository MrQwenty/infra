#!/bin/bash

# Test script for WhatsApp verification system
# This script tests the WhatsApp Meta Business API integration

echo "🚀 Testing WhatsApp Meta Business API Integration"
echo "================================================"

# Configuration
API_BASE_URL="http://localhost:3231"
WHATSAPP_API_URL="https://graph.facebook.com/v19.0/676124925591256/messages"
ACCESS_TOKEN="EAA6YlMvSeosBPM6NEv0SDKPu5IrTuAkLqL3jsQPglG181ZBD2bLy9P0TEtFJZBr064A5PFSc3fZAuZCMJeUKCkYbs2CN0vJkkuRLPJXqbaP8bpuX3ZC3PtX1yh7ZCexyjgTSjTy6PVUujRQ9cydJ4XV8ZBxYGRojZArolTo14YBnCgoJaf2VUdZBy1zOsQ4AyF1JrKD3gB64w8pBvhSTN1sDPyVjsSsoOiM25FyerpVdsqBVHaAZDZD"
PHONE_NUMBER_ID="676124925591256"

# Test phone number (use your own WhatsApp number for testing)
TEST_PHONE_NUMBER="+1234567890"  # Replace with your actual WhatsApp number

echo "📱 Testing direct WhatsApp API call..."

# Test direct WhatsApp API call
curl -X POST "$WHATSAPP_API_URL" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "'$TEST_PHONE_NUMBER'",
    "type": "template",
    "template": {
      "name": "hello_world",
      "language": {
        "code": "en_US"
      },
      "components": [
        {
          "type": "body",
          "parameters": [
            {
              "type": "text",
              "text": "Your InfluenzaNet verification code is: 123456. This code will expire in 10 minutes."
            }
          ]
        }
      ]
    }
  }' | jq '.'

echo ""
echo "✅ Direct API test completed. Check your WhatsApp for the message."
echo ""

# Test the application endpoints (requires authentication)
echo "🔧 Testing application endpoints..."
echo "Note: You need to be logged in to test these endpoints"
echo ""

echo "1. Test add phone endpoint:"
echo "POST $API_BASE_URL/v1/user/contact/add-phone"
echo "Body: {\"newPhone\": \"$TEST_PHONE_NUMBER\", \"verificationMethod\": \"whatsapp\"}"
echo ""

echo "2. Test verify phone endpoint:"
echo "POST $API_BASE_URL/v1/user/contact/verify-phone"
echo "Body: {\"token\": \"verification_token\", \"code\": \"123456\"}"
echo ""

echo "📋 Environment Check:"
echo "- WhatsApp API URL: $WHATSAPP_API_URL"
echo "- Phone Number ID: $PHONE_NUMBER_ID"
echo "- Access Token: ${ACCESS_TOKEN:0:20}..."
echo ""

echo "🔍 Troubleshooting Tips:"
echo "1. Make sure your WhatsApp Business Account is verified"
echo "2. Ensure the 'hello_world' template is approved"
echo "3. Check that the phone number is registered with WhatsApp"
echo "4. Verify the access token hasn't expired"
echo ""

echo "✨ Test completed! Check the logs for any errors."