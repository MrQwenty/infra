# WhatsApp Notification Retry System

## Overview

This document describes the enhanced automated retry mechanism for failed WhatsApp notification deliveries implemented in the InfluenzaNet messaging service.

## Features

### Exponential Backoff Strategy
- **Base Delay**: 30 seconds for normal messages, 15 seconds for verification codes
- **Maximum Delay**: 1 hour (3600 seconds) to prevent indefinite delays
- **Jitter**: 10% randomization to prevent thundering herd problems
- **Error-Specific Strategies**: Different backoff multipliers based on error type

### Error Categorization

The system categorizes errors into four types:

1. **Rate Limit Errors** (`rate_limit`)
   - Triggers: "rate limit", "429" HTTP status
   - Strategy: Aggressive exponential backoff (2^retry_count * 2)
   - Retryable: Yes

2. **Network Errors** (`network_error`)
   - Triggers: "network", "timeout"
   - Strategy: Moderate backoff (1.5^retry_count)
   - Retryable: Yes

3. **API Errors** (`api_error`)
   - Triggers: General API failures
   - Strategy: Standard exponential backoff (2^retry_count)
   - Retryable: Yes

4. **Invalid Number** (`invalid_number`)
   - Triggers: "invalid" + "number"
   - Strategy: No retry (immediate deletion)
   - Retryable: No

### Database Schema Enhancements

The `OutgoingWhatsApp` collection now includes:

```go
type OutgoingWhatsApp struct {
    ID               primitive.ObjectID `bson:"_id,omitempty"`
    MessageType      string             `bson:"messageType"`
    To               string             `bson:"to"`
    Content          string             `bson:"content"`
    AddedAt          int64              `bson:"addedAt"`
    HighPrio         bool               `bson:"highPrio"`
    LastSendAttempt  int64              `bson:"lastSendAttempt"`
    RetryCount       int                `bson:"retryCount"`
    MaxRetries       int                `bson:"maxRetries"`
    NextRetryAt      int64              `bson:"nextRetryAt"`     // New field
    BaseDelaySeconds int                `bson:"baseDelaySeconds"` // New field
    LastErrorType    string             `bson:"lastErrorType"`   // New field
}
```

## Configuration

### Environment Variables

- `WHATSAPP_MESSAGE_SCHEDULER_INTERVAL`: Scheduler check interval (default: 10 seconds)
- `WHATSAPP_CLIENT_SERVICE_LISTEN_PORT`: WhatsApp client service port (default: 5007)

### Message Type Configuration

- **Normal Messages**: MaxRetries=5, BaseDelay=30s
- **Verification Messages**: MaxRetries=3, BaseDelay=15s (faster retry for time-sensitive codes)

## User Channel Preferences

The system respects user notification preferences stored in `ContactPreferences`:

```go
type ContactPreferences struct {
    SubscribedToWhatsApp  bool     `bson:"subscribedToWhatsApp"`
    WhatsAppNumber        string   `bson:"whatsAppNumber"`
    PreferredChannels     []string `bson:"preferredChannels"`
}
```

### Channel Selection Logic

1. Check user's `PreferredChannels` array
2. If empty, default to `["email"]`
3. For each preferred channel:
   - **Email**: Queue email notification
   - **WhatsApp**: Check subscription status and phone number availability
4. Success if any channel succeeds

## Monitoring and Debugging

### Log Levels

- **Info**: Retry scheduling with delay calculations
- **Debug**: Successful message delivery
- **Warning**: Max retries exceeded, non-retryable errors
- **Error**: Database operations, client failures

### Key Log Messages

```
WhatsApp message {id} scheduled for retry in {delay} seconds (attempt {current}/{max})
WhatsApp message {id} has non-retryable error, deleting
WhatsApp message {id} exceeded max retries, deleting
```

## Performance Considerations

### Batch Processing
- High priority messages processed first
- Batch size: 50 messages per instance per cycle
- Separate queues for high/normal priority

### Database Optimization
- Index on `nextRetryAt` for efficient retry scheduling
- Index on `highPrio` for priority-based fetching
- Atomic operations for retry count updates

### Memory Usage
- Exponential backoff prevents memory buildup from failed messages
- Maximum delay cap prevents indefinite queue growth
- Non-retryable errors immediately removed

## Migration Notes

### Backward Compatibility
- Existing messages without `nextRetryAt` are processed immediately
- Default values applied for missing `baseDelaySeconds`
- Legacy retry logic gracefully handled

### Database Migration
No explicit migration required. New fields are optional and have sensible defaults.

## Testing

### Unit Tests
```bash
go test ./messaging-service/pkg/retry/...
```

### Integration Testing
1. Mock WhatsApp API with different error responses
2. Verify exponential backoff timing
3. Test error categorization logic
4. Confirm non-retryable errors are deleted

### Load Testing
- Test with high message volumes
- Verify scheduler performance under load
- Monitor database query performance

## Troubleshooting

### Common Issues

1. **Messages not retrying**
   - Check `nextRetryAt` timestamp
   - Verify scheduler is running
   - Check error categorization

2. **Excessive retries**
   - Review `MaxRetries` configuration
   - Check error classification logic
   - Monitor exponential backoff calculations

3. **Performance degradation**
   - Review batch sizes
   - Check database indexes
   - Monitor scheduler interval

### Debug Commands

```bash
# Check outgoing messages
db.outgoing_whatsapp.find({nextRetryAt: {$lte: new Date().getTime()/1000}})

# Monitor retry patterns
db.outgoing_whatsapp.aggregate([
  {$group: {_id: "$lastErrorType", count: {$sum: 1}}}
])
```

## Future Enhancements

- Dead letter queue for permanently failed messages
- Retry pattern analytics and optimization
- Dynamic retry strategy based on error patterns
- Circuit breaker pattern for failing WhatsApp providers
