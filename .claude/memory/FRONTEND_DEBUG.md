# Frontend Debugging Guide

This document explains the debugging mechanism implemented in the ExpertDB frontend application to assist with development and troubleshooting.

## Debug Mode

The application includes a configurable debug mode that controls logging verbosity and helps with API troubleshooting.

### Configuration

Debug mode is controlled by the `VITE_DEBUG_MODE` environment variable:

```
# In .env, .env.development, or .env.production
VITE_DEBUG_MODE=true  # Enable debug mode
VITE_DEBUG_MODE=false # Disable debug mode
```

When not specified, debug mode defaults to disabled.

### Usage

Debug mode primarily affects API-related logging in the `api.ts` service:

```javascript
// Check if we're in debug mode
const isDebugMode = import.meta.env.VITE_DEBUG_MODE === 'true';
```

### Logging Behavior

#### When Debug Mode is Enabled

1. **Request Logging**:
   - Logs URL, method, and parameters for all API requests
   - Example: `Making request: GET /experts {name: "John"}`

2. **Response Logging**:
   - Logs complete response data from successful API calls
   - Example: `Response from /experts: [{id: "1", name: "John Doe", ...}]`

3. **Error Logging**:
   - Logs detailed error information including:
     - Error status code
     - Error response data
     - Error message and stack trace
   - Example: `Error 404 for /experts/999: {"success": false, "message": "Expert not found"}`

4. **Authentication**:
   - Logs authentication attempts (with sensitive data omitted)
   - Example: `Sending login request to: /api/auth/login`

#### When Debug Mode is Disabled

1. **Minimal Logging**:
   - Only critical errors are logged
   - No sensitive data is exposed in logs
   - No request/response payload logging

2. **Production Behavior**:
   - User-friendly error messages are still displayed in the UI
   - Backend error messages are used when available
   - Generic messages are used for network or unexpected errors

## Debugging API Issues

When encountering API-related issues:

1. Enable debug mode in your local environment
2. Check browser console logs for detailed request/response information
3. Look for specific error status codes and messages
4. Examine network requests in browser developer tools

## Common Error Status Codes

The debug mode helps identify specific HTTP status errors:

- **400 (Bad Request)**: Invalid data sent to API
- **401 (Unauthorized)**: Authentication required or invalid token
- **403 (Forbidden)**: Insufficient permissions for the operation
- **404 (Not Found)**: Resource does not exist
- **500 (Server Error)**: Backend error processing the request

## Adding Debug Logs

When adding new features or troubleshooting, you can add your own debug logs:

```javascript
if (isDebugMode) {
  console.log('Debug information:', someValue);
}
```

This ensures logs only appear in development environments and not in production.

## Security Considerations

- Never log sensitive information like passwords, tokens, or personal data, even in debug mode
- For sensitive operations, log that the operation occurred but not the specific data
- Use the debug mode condition to prevent exposing sensitive information

## Best Practices

1. Keep debug logs informative but concise
2. Group related debug information together
3. Clear unnecessary debug logs before deploying to production
4. Consider adding environment-specific logging levels for more granular control