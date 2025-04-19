# Chirpy API

# Installation
`go install github.com/sambakker4/chirpy`

## Create and setup database
Make sure you have `goose`, `postgres`, and `git` installed
> Run `git clone github.com/sambakker4/chirpy`
> Create a local database called `chirpy` with and save the url in the `.env` file at the base of the repo as `DB_URL`
> Run `goose postgres your_database_url up` in the `sql/schema` directory

## Error
Any errors will be returned in the following format
```json
{
    "error":"what happened"
}
```

## User Resources

### POST /api/users

Create a new user profile
> Request
```json
{
    "email":"something@example.com",
    "password":"password123"
}
```

> Response
```json
{
    "id":"c3eded59-9421-4d38-9f9a-291e89ce993a",
    "created_at":"2025-04-19T14:48:11.338278Z",
    "updated_at":"2025-04-19T14:48:11.338278Z",
    "email":"something@example.com",
    "is_chirpy_red":false
}
```

### PUT /api/user

Update user info (for refresh_token or token see */api/login*)
> Request
Must have a header `Authorization` that is set to a token in the format `Bearer <token>`
```json
{
    "email":"something@example.com",
    "password":"password123"
}
```

> Response
```json
{
  "id": "c3eded59-9421-4d38-9f9a-291e89ce993a",
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "0001-01-01T00:00:00Z",
  "email": "something@example.com",
  "token": "",
  "refresh_token": "",
  "is_chirpy_red": false
}
```

### POST /api/login

Login user
> Request
```json
{
    "email":"something@example.com",
    "password":"password123"
}
```

> Response
```json
{
    "id":"c3eded59-9421-4d38-9f9a-291e89ce993a",
    "created_at":"2025-04-19T14:48:11.338278Z",
    "updated_at":"2025-04-19T14:48:11.338278Z",
    "email":"something@example.com",
    "token":"eyjhbgcioijiuzi1niisinr5cci6ikpxvcj9.eyjpc3mioijjaglychkilcjzdwiioijjm2vkzwq1os05ndixltrkmzgtowy5ys0yotflodljztk5m2eilcjlehaioje3nduwoti1odesimlhdci6mtc0nta4odk4mx0.hjt_shtftabnb_5e397i1it75aky_hh_c1t4qdai6es",
    "refresh_token":"5148691ae5a0a38e3b330edf492fb4b3bf1032b14d9f288ba32c0cd8e6930ef8",
    "is_chirpy_red":false
}
```
**Tokens are for authentication and are used to access resources as a certain user** 
**Refresh tokens are used get another token at the endpoint */api/refresh***

### POST /api/refresh
Must have a header `Authorization` that is set to a token in the format `Bearer <refresh_token>`

> Request
/api/refresh

> Response
```json
{
    "token":"eyjhbgcioijiuzi1niisinr5cci6ikpxvcj9.eyjpc3mioijjaglychkilcjzdwiioijjm2vkzwq1os05ndixltrkmzgtowy5ys0yotflodljztk5m2eilcjlehaioje3nduwoti1odesimlhdci6mtc0nta4odk4mx0.hjt_shtftabnb_5e397i1it75aky_hh_c1t4qdai6es"
}
```

### Get /api/chirps
Returns an array of all chrips
*Has an optional query parameter of sort that can be set to `asc` or `desc`* (`asc` is the default)

> Response
```json
  {
    "id": "50bac20a-25bc-4442-9cbb-177ce6b8c4a0",
    "created_at": "2025-04-19T12:50:06.992009Z",
    "updated_at": "2025-04-19T12:50:06.992009Z",
    "body": "Gale!",
    "user_id": "9323084d-f748-4666-a99c-a1804de93066"
  },
  {
    "id": "96163344-d5a6-4a75-afc5-6caea49f0ce5",
    "created_at": "2025-04-19T12:50:06.993199Z",
    "updated_at": "2025-04-19T12:50:06.993199Z",
    "body": "Cmon Pinkman",
    "user_id": "9323084d-f748-4666-a99c-a1804de93066"
  },
]
```
### POST /api/chirps

Update user info (for refresh_token or token see */api/login*)
> Request
Must have a header `Authorization` that is set to a token in the format `Bearer <token>`

```json
{
    "email":"something@example.com",
    "password":"password123",
    "body":"message"
}
```

> Response
```json
{
  "id": "ed01cd2d-f739-48cd-8d81-ed54bfdd2c8d",
  "created_at": "2025-04-19T17:54:16.499063Z",
  "updated_at": "2025-04-19T17:54:16.499063Z",
  "body": "message",
  "user_id": "c3eded59-9421-4d38-9f9a-291e89ce993a"
}
```

### GET /api/chirps/{chirpID}

Get the chrip with the matching chrip id
> Request
/api/chirps/ed01cd2d-f739-48cd-8d81-ed54bfdd2c8d 

> Response
```json
{
  "id": "ed01cd2d-f739-48cd-8d81-ed54bfdd2c8d",
  "created_at": "2025-04-19T17:54:16.499063Z",
  "updated_at": "2025-04-19T17:54:16.499063Z",
  "body": "message",
  "user_id": "c3eded59-9421-4d38-9f9a-291e89ce993a"
}
```

### DELETE /api/chirps/{chirpID}
Delete chirp (for refresh_token or token see */api/login*)
> Request
Must have a header `Authorization` that is set to a token in the format `Bearer <token>`
/api/chirps/ed01cd2d-f739-48cd-8d81-ed54bfdd2c8d

> Response
Nothing

### POST /api/revoke
Revokes refresh token
> Request
Must have a header `Authorization` that is set to a token in the format `Bearer <refresh_token>`
/api/revoke

> Reponse
Nothing

### POST /api/polka/webhook
This example upgrade user to Chirpy Red
Must have a header `Authorization` that is set to a token in the format `ApiKey <api_key>`

> Request
```json
{
    "event":"user.upgrade",
    "data": {
        "user_id": "c3eded59-9421-4d38-9f9a-291e89ce993a",
    }
}
```

> Response
Nothing
