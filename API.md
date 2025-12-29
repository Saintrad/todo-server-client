todo v1

1. Overview

Purpose: Provide a Todo backend service consumed by a CLI client.

Transport: HTTP

Format: JSON (UTF-8)

Base URLs:

v1: http://localhost:{port}/v1

v2: http://localhost:{port}/v2

Content-Type: application/json; charset=utf-8

Time format: RFC3339 (e.g., 2025-12-29T12:34:56Z)

ID format: Server-generated integer (monotonic within a datastore)

1.1 Conventions

JSON field naming: snake_case

Booleans are explicit (true/false).

Unknown fields in request bodies should be rejected with 400 VALIDATION_ERROR (recommended).

2. Resource Models
2.1 Todo (Response DTO)

Fields returned by the API.

id (int, required) — server-generated

title (string, required) — non-empty, recommended max length 200

category (string, optional)

is_done (bool, required)

due_at (string|null, optional) — RFC3339, nullable

created_at (string, required) — RFC3339

updated_at (string, required) — RFC3339

Notes / invariants

title must not be blank.

created_at set by server on creation.

updated_at set by server on creation and on any update.

is_done defaults to false on create.

due_at is optional; omit or set to null to represent “no due date.”

3. Error Model (Standard)

All errors return JSON with this shape:

{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "title is required",
    "details": [
      { "field": "title", "issue": "must not be empty" }
    ]
  }
}

3.1 Error Codes

VALIDATION_ERROR → 400

NOT_FOUND → 404

UNAUTHORIZED → 401 (v2 only)

FORBIDDEN → 403 (v2 only; optional if you add roles/ownership checks beyond scoping)

CONFLICT → 409 (e.g., username already exists)

INTERNAL → 500

Part A — v1 API (No Auth, Single-User)
A1. Endpoints Summary

POST /v1/todos — create

GET /v1/todos — list

GET /v1/todos/{id} — get one

PATCH /v1/todos/{id} — update (partial)

DELETE /v1/todos/{id} — delete

A1.1 Create Todo

Method: POST

Path: /v1/todos

Request body:

{
  "title": "Buy milk",
  "category": "errands",
  "due_at": "2026-01-05T10:00:00Z"
}


Responses:

201 Created → body: Todo

400 Bad Request → body: Error

A1.2 List Todos

Method: GET

Path: /v1/todos

Query params (all optional):

is_done: true|false

q: string (search in title/category)

sort: created_at|due_at|updated_at (default: created_at)

order: asc|desc (default: desc)

limit: int (default: 50, max: 200)

offset: int (default: 0)

Responses:

200 OK → body:

{
  "items": [ /* Todo[] */ ],
  "total": 123,
  "limit": 50,
  "offset": 0
}

A1.3 Get Todo

Method: GET

Path: /v1/todos/{id}

Responses:

200 OK → body: Todo

404 Not Found → body: Error

A1.4 Update Todo (Partial)

Method: PATCH

Path: /v1/todos/{id}

PATCH semantics

Omitted field → unchanged

null → clear value (only for nullable fields such as due_at and category)

Request body (any subset allowed):

{
  "title": "Buy milk and eggs",
  "category": null,
  "is_done": true,
  "due_at": null
}


Responses:

200 OK → body: Todo

400 Bad Request → body: Error

404 Not Found → body: Error

A1.5 Delete Todo

Method: DELETE

Path: /v1/todos/{id}

Responses:

204 No Content

404 Not Found → body: Error