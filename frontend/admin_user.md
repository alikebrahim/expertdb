WAP Error: Request failed with status code 500

SHOW_ERROR_LISTED_FRAMES

[Anonymous] then
Promise
createUser
orSubmit
e userForm.tsx:64

SHOW_ERROR_LISTED_FRAMES

[Anonymous] then
Promise
createUser
orSubmit
e userForm.tsx:64

SHOW_ERROR_LISTED_FRAMES

Error: Request failed with status code 500; name: "AxiosError"; code: "ERR_BAD_RESPONSE"; config: {..., method: "POST", url: "http://localhost:8808/api/users", ...}; adapter: "xhr"; transformerRequest: [Array(1)]; transformerResponse: [Array(1)]; ...

Response: {
  "error": {
    "status": 500,
    "statusText": "Internal Server Error",
    "headers": {
      "axiosHeaders": {
        "config": {...},
        "upload": {}
      }
    },
    "data": {
      "error": "Failed to create user: failed to check for existing email: no such table: users"
    }
  }
}

SHOW_ERROR_LISTED_FRAMES

awaitIn request
createUser
orSubmit
e userForm.tsx:64
