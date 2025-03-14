/admin - User Management tab:
No Console Errors: Which is weird!

/admin - Expert Requests tab:
## Browser Console Errors

### Error 1: HTTP Request Error
- **Type**: HTTP Request Error
- **Method**: `GET`
- **URL**: `http://localhost:8008/api/expert-requests`
- **Status Code**: `500 (Internal Server Error)`
- **File**: `api.ts:103`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `AdminPage.tsx:00`
  - `(anonymous)` @ `AdminPage.tsx:88`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 2: API Error (Uncaught in Promise)
- **Type**: API Error (Uncaught in Promise)
- **Message**: `Failed with stahttp://localhost:5173/requeststus code 500`
- **File**: `api.ts:39`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `AdminPage.tsx:00`
  - `(anonymous)` @ `AdminPage.tsx:88`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 3: Response Data Error
- **Type**: Response Data Error
- **Message**: `(anonymous) e [error: 'Failed to retrieve expert requests: failed to query.expert_requests: no such column: rejection_reason']`
- **File**: `api.ts:48`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `AdminPage.tsx:00`
  - `(anonymous)` @ `AdminPage.tsx:88`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 4: Status 500 Error
- **Type**: Status 500 Error
- **Context**: `(anonymous)`
- **File**: `api.ts:41`
- **Stack Trace**:
  - `Promise.then`
  - `request`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `AdminPage.tsx:00`
  - `(anonymous)` @ `AdminPage.tsx:88`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 5: Request Failure (AxiosError)
- **Type**: Request Failure (AxiosError)
- **Message**: `Request failed with status code 500`
- **Error Name**: `AxiosError`
- **Error Code**: `ERR_BAD_RESPONSE`
- **Configuration**: `(...)`
- **Request Type**: `XMLHttpRequest`
- **File**: `api.ts:07`
- **Stack Trace**:
  - `await in getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `AdminPage.tsx:00`
  - `(anonymous)` @ `AdminPage.tsx:88`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 6: Error Response Data
- **Type**: Error Response Data
- **Message**: `error: 'Failed to retrieve expert requests: failed to query.expert_requests: no such column: rejection_reason'`
- **File**: `api.ts:72`
- **Stack Trace**:
  - `await in getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `AdminPage.tsx:00`
  - `(anonymous)` @ `AdminPage.tsx:88`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`


/search:
When user logs in, after redirect to /search
## Browser Console Errors (New Set)

### Error 1: TypeError (Cannot Read Properties of Null)
- **Type**: TypeError
- **Message**: `Cannot read properties of null (reading 'success')`
- **File**: `SearchPage.tsx:46`
- **Stack Trace**:
  - `at fetchExperts` @ `SearchPage.tsx:46`
  - `fetchExperts` @ `SearchPage.tsx:34`
  - `(anonymous)` @ `SearchPage.tsx:01`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`
- **Response**: `Response from /api/experts: null`

---

### Error 2: TypeError (Cannot Read Properties of Null)
- **Type**: TypeError
- **Message**: `Cannot read properties of null (reading 'success')`
- **File**: `SearchPage.tsx:46`
- **Stack Trace**:
  - `await in fetchExperts` @ `SearchPage.tsx:34`
  - `fetchExperts` @ `SearchPage.tsx:34`
  - `(anonymous)` @ `SearchPage.tsx:01`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`
- **Response**: `Response from /api/experts: null`


/requests:
when user vists /requests through Expert Requests button on side panel:
## Browser Console Errors

### Error 1: HTTP Request Error
- **Type**: HTTP Request Error
- **Method**: `GET`
- **URL**: `http://localhost:8008/api/expert-requests`
- **Status Code**: `500 (Internal Server Error)`
- **File**: `api.ts:103`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 2: API Error (Uncaught in Promise) (Screenshot 1)
- **Type**: API Error (Uncaught in Promise)
- **Message**: `Failed with status code 500`
- **File**: `api.ts:39`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 3: Response Data Error (Screenshot 1)
- **Type**: Response Data Error
- **Message**: `[error: 'Failed to retrieve expert requests: failed to query.expert_requests: no such column: rejection_reason']`
- **File**: `api.ts:48`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 4: Status 500 Error (Screenshot 1)
- **Type**: Status 500 Error
- **Context**: `(anonymous)`
- **File**: `api.ts:41`
- **Stack Trace**:
  - `Promise.then`
  - `request`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 5: Request Failure (AxiosError) (Screenshot 1)
- **Type**: Request Failure (AxiosError)
- **Message**: `Request failed with status code 500`
- **Error Name**: `AxiosError`
- **Error Code**: `ERR_BAD_RESPONSE`
- **Configuration**: `(...)`
- **Request Type**: `XMLHttpRequest`
- **File**: `api.ts:07`
- **Stack Trace**:
  - `await in getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 6: Error Response Data (Screenshot 1)
- **Type**: Error Response Data
- **Message**: `[error: 'Failed to retrieve expert requests: failed to query.expert_requests: no such column: rejection_reason']`
- **File**: `api.ts:72`
- **Stack Trace**:
  - `await in getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 7: HTTP Request Error 
- **Type**: HTTP Request Error
- **Method**: `GET`
- **URL**: `http://localhost:8008/api/expert-requests?userId=2`
- **Status Code**: `500 (Internal Server Error)`
- **File**: `api.ts:103`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 8: API Error (Uncaught in Promise) 
- **Type**: API Error (Uncaught in Promise)
- **Message**: `Failed with status code 500`
- **File**: `api.ts:39`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 9: Response Data Error 
- **Type**: Response Data Error
- **Message**: `[error: 'Failed to retrieve expert requests: failed to query.expert_requests: no such column: rejection_reason']`
- **File**: `api.ts:48`
- **Stack Trace**:
  - `Promise.then`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 10: Status 500 Error 
- **Type**: Status 500 Error
- **Context**: `(anonymous)`
- **File**: `api.ts:41`
- **Stack Trace**:
  - `Promise.then`
  - `request`
  - `getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 11: Request Failure (AxiosError) 
- **Type**: Request Failure (AxiosError)
- **Message**: `Request failed with status code 500`
- **Error Name**: `AxiosError`
- **Error Code**: `ERR_BAD_RESPONSE`
- **Configuration**: `(...)`
- **Request Type**: `XMLHttpRequest`
- **File**: `api.ts:07`
- **Stack Trace**:
  - `await in getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 12: Error Response Data 
- **Type**: Error Response Data
- **Message**: `[error: 'Failed to retrieve expert requests: failed to query.expert_requests: no such column: rejection_reason']`
- **File**: `api.ts:72`
- **Stack Trace**:
  - `await in getExpertRequests` @ `api.ts:102`
  - `fetchRequests` @ `ExpertRequestPage.tsx:27`
  - `(anonymous)` @ `ExpertRequestPage.tsx:44`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`


/requests:
when user tries to create a new Expert Request:
## Browser Console Errors (New Set)

### Error 1: HTTP Request Error
- **Type**: HTTP Request Error
- **Method**: `POST`
- **URL**: `http://localhost:8008/api/expert-requests`
- **Status Code**: `400 (Bad Request)`
- **File**: `api.ts:103`
- **Stack Trace**:
  - `Promise.then`
  - `createExpertRequest` @ `api.ts:173`
  - `onSubmit` @ `ExpertRequestForm.tsx:81`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 2: API Error (Uncaught in Promise)
- **Type**: API Error (Uncaught in Promise)
- **Message**: `Request failed with status code 400`
- **File**: `api.ts:39`
- **Stack Trace**:
  - `Promise.then`
  - `createExpertRequest` @ `api.ts:173`
  - `onSubmit` @ `ExpertRequestForm.tsx:81`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 3: Response Data Error
- **Type**: Response Data Error
- **Message**: `(anonymous) e [error: 'Invalid request payload']`
- **File**: `api.ts:48`
- **Stack Trace**:
  - `Promise.then`
  - `createExpertRequest` @ `api.ts:173`
  - `onSubmit` @ `ExpertRequestForm.tsx:81`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 4: Status 400 Error
- **Type**: Status 400 Error
- **Context**: `(anonymous)`
- **File**: `api.ts:41`
- **Stack Trace**:
  - `Promise.then`
  - `createExpertRequest` @ `api.ts:173`
  - `onSubmit` @ `ExpertRequestForm.tsx:81`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`

---

### Error 5: Request Failure (AxiosError)
- **Type**: Request Failure (AxiosError)
- **Message**: `Request failed with status code 400`
- **Error Name**: `AxiosError`
- **Error Code**: `ERR_BAD_REQUEST`
- **Configuration**: `(...)`
- **Request Type**: `XMLHttpRequest`
- **File**: `api.ts:07`
- **Stack Trace**:
  - `await in createExpertRequest` @ `api.ts:173`
  - `onSubmit` @ `ExpertRequestForm.tsx:81`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`
- **Details**:
  - `config`: `{ url: 'http://localhost:8008/api/expert-requests', method: 'POST', transformRequest: Array[1], transformResponse: Array[1], timeout: 0, adapter: Array[3], ... }`
  - `message`: `Request failed with status code 400`
  - `name`: `AxiosError`
  - `response`: `{ status: 400, statusText: 'Bad Request', headers: AxiosHeaders, config: {...}, request: XMLHttpRequest, ... }`
  - `stack`: `AxiosError: Request failed with status code 400\n    at settle (.../node_modules/axios/...)\n    at XMLHttpRequest.onloadend (.../node_modules/vite/...)\n    at XMLHttpRequest.dispatchEvent (...)\n    at XMLHttpRequest.setReadyState (...)\n    at XMLHttpRequest.__didCompleteResponse__ (...)\n    at ...`

---

### Error 6: Error Response Data
- **Type**: Error Response Data
- **Message**: `[error: 'Invalid request payload']`
- **File**: `api.ts:72`
- **Stack Trace**:
  - `await in createExpertRequest` @ `api.ts:173`
  - `onSubmit` @ `ExpertRequestForm.tsx:81`
- **Additional Action**: `SHOW _IGNORABLE_LISTED_FRAMES`
