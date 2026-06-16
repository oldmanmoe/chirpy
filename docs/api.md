# USERS

## POST /api/users
Creates a new user account.

* **Auth Required:** No

### Request Body (JSON)
* `email` (string, required): Unique email address.
* `password` (string, required): Plaintext password.

### Responses
* **201 Created**: Returns the user object (minus the password).
* **400 Bad Request**: Missing fields or email already taken.

---

## POST /api/login
Authenticates a user and returns a token.

* **Auth Required:** No

### Request Body (JSON)
* `email` (string, required)
* `password` (string, required)

### Responses
* **200 OK**: Returns a JSON payload containing the access token.
* **401 Unauthorized**: Invalid email or password.

---

## POST /api/polka/webhooks
Handles incoming upgrade events from a third-party payment processor.

* **Auth Required:** Yes (API Key in `Authorization` header)

### Request Body (JSON)
* `event` (string, required): The type of event (e.g., "user.upgraded").
* `data` (object, required): Contains the payload.
  * `user_id` (string/UUID, required): The ID of the user to upgrade.

### Responses
* **204 No Content**: Successfully processed and user status updated.
* **401 Unauthorized**: Invalid or missing API key.
* **404 Not Found**: The `user_id` provided doesn't exist.

---

## POST /api/users
Updates

* **Auth Required:** Yes (API Key in `Authorization` header)

### Request Body (JSON)
* `email` (string, required)
* `password` (string, required)

### Responses
* **200 OK**: Returns a JSON payload containing the access token.
* **401 Unauthorized**: Invalid or missing API Key.
* **500 Internal Server Error**: Database connection failure or unhandled server error. Returns a generic error message.

---

# CHIRPS

## POST /api/chirps

** **Auth Requires:** No

### Request Body (JSON)
* `body` (string, required)
* `user_id` (string/UUID, required)

### Responses

* **201 Created**: Retruns a Json payload containing details of a specific chirp by its unique ID.

```json
{
  "id": "aa5a315d61ae9438b18d",
  "created_at": "2024-02-13T12:00:00Z",
  "updated_at": "2024-02-13T12:00:00Z",
  "body": "example",
  "user_id": null
}
```
* **400 Bad Request**: The server couldn't process the request because of invalid data.
* **401 Unauthorized**: Invalid or missing API.
* **500 Internal Server Error**: Database connection failure or unhandled server error. Returns a generic error message.