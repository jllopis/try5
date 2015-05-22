TRY5
====

Version: 0.0.1

Try5 is a new intent to develop an easy and basic authentication system. The functionalyty is exposed through REST Web Services.

It is an opinionated library with absolutely no garanties. Use at your own risk.

You need:

- Go v1.4.2
- PostgreSQL database v9.4 (optional)
- Not much really...

Proposed features:

- HMAC Authentication
- JWT Authentication
- Plain user / password authentication
- User management
- Key/Token managemnt
- RBAC Authorization
- REST endpoint
- RPC endpoint

Getting Started
---------------

### Install

~~~
go get github.com/jllopis/try5
~~~

Specification
-------------

The server expect to receive the data in the body of the request for `POST` and `PUT` verbs using type `application/json; charset=UTF-8`. The data provided as _url vars_ or via `application/x-www-form-urlencoded` **will not be accepted**.

* `GET` request to `/api/v1/accounts`

	````
	$ curl -ki https://b2d:9000/api/v1/accounts
	HTTP/1.1 200 OK
	Content-Type: application/json; charset=UTF-8
	Date: Fri, 22 May 2015 11:23:57 GMT
	Content-Length: 715
	
	[
	  {
	    "uid": "447fb74b-114c-46c9-aee4-292998d845bb",
	    "email": "tu5@test.com",
	    "name": "test user 5",
	    "password": "$2a$10$h.tOLA00ZLKux5IWa6l4JOcxpa2SdCTS56w2p7I31JZw2lZM3ivIK",
	    "active": true,
	    "gravatar": null,
	    "created": "2015-05-22T11:15:05.840968723Z",
	    "updated": "2015-05-22T11:15:05.840968723Z",
	    "deleted": null
	  },
	  {
	    "uid": "60b51e16-fe83-4ac2-853c-7cbc1f250a09",
	    "email": "tu4@test.com",
	    "name": "test user 4",
	    "password": "$2a$10$KfX8Sf13JsQWSelCxDDFluy0iHMaNmVYzGrZTHux3NYKS7476Qp0y",
	    "active": true,
	    "gravatar": null,
	    "created": "2015-05-22T11:22:32.145080999Z",
	    "updated": "2015-05-22T11:22:32.145080999Z",
	    "deleted": null
	  }
	]
	````

* `GET` request to `/api/v1/accounts/7ecee355-537b-492c-ab23-6a41219959d1`

	````
	$ curl -ki https://b2d:9000/api/v1/accounts/7ecee355-537b-492c-ab23-6a41219959d1
	HTTP/1.1 200 OK
	Content-Type: application/json; charset=UTF-8
	Date: Fri, 22 May 2015 10:57:40 GMT
	Content-Length: 359
	
	[
	  {
	    "uid": "7ecee355-537b-492c-ab23-6a41219959d1",
	    "email": "tu5@test.com",
	    "name": "test user 5",
	    "password": "$2a$10$651aQ.FoQQo3t2BBhI3WX.rVSKbPu4ks8feMG1zi5dBvysraY7nYm",
	    "active": true,
	    "gravatar": null,
	    "created": "2015-05-22T10:01:56.160527217Z",
	    "updated": "2015-05-22T10:01:56.160527217Z",
	    "deleted": null
	  }
	]
	````

* `POST` request to `/api/v1/accounts`

	````
	$ curl -ki https://b2d:9000/api/v1/accounts -X POST -d '{"email":"tu4@test.com","name":"test user 4","password":"1234","active":true}'
	HTTP/1.1 200 OK
	Content-Type: application/json; charset=UTF-8
	Date: Fri, 22 May 2015 11:22:32 GMT
	Content-Length: 333
	
	{
	  "uid": "60b51e16-fe83-4ac2-853c-7cbc1f250a09",
	  "email": "tu4@test.com",
	  "name": "test user 4",
	  "password": "$2a$10$KfX8Sf13JsQWSelCxDDFluy0iHMaNmVYzGrZTHux3NYKS7476Qp0y",
	  "active": true,
	  "gravatar": null,
	  "created": "2015-05-22T11:22:32.145080999Z",
	  "updated": "2015-05-22T11:22:32.145080999Z",
	  "deleted": null
	}
	````

* `PUT` request to `/api/v1/accounts/`

	````
	$ curl -ki https://b2d:9000/api/v1/accounts/e557e74a-cb35-4039-b4e5-f9c6ca777c5b -X PUT -d '{"name": "Test User 4","email":"newtu4@test4.com","password":"1234","active":true}'
	HTTP/1.1 200 OK
	Content-Type: application/json; charset=UTF-8
	Date: Fri, 22 May 2015 11:48:41 GMT
	Content-Length: 309
	
	{
	  "uid": "e557e74a-cb35-4039-b4e5-f9c6ca777c5b",
	  "email": "newtu4@test4.com",
	  "name": "Test User 4",
	  "password": "$2a$10$MWCvQXeCw0D1jXYQUMGCJuAFsTPzTvuYYVE2/1pEhu/.LQHqmqsPu",
	  "active": true,
	  "gravatar": null,
	  "created": "2015-05-22T11:22:32.145080999Z",
	  "updated": "2015-05-22T11:48:41.567466863Z",
	  "deleted": null
	}
	````

* `DELETE` request to `/api/v1/accounts/802aa9ef-b00e-4204-9b75-4dbb82d20643`

	````
	$ curl -ki https://localhost:9000/api/v1/accounts/802aa9ef-b00e-4204-9b75-4dbb82d20643 -X DELETE
	HTTP/1.1 200 OK
	Content-Type: application/json; charset=UTF-8
	Date: Fri, 22 May 2015 15:56:49 GMT
	Content-Length: 164
	
	{
	  "status": "ok",
	  "action": "delete",
	  "info": "802aa9ef-b00e-4204-9b75-4dbb82d20643",
	  "table": "accounts",
	  "id": "802aa9ef-b00e-4204-9b75-4dbb82d20643"
	}
	````

* `POST` request to `/api/v1/authenticate`

	````
	$ curl -ki https://localhost:9000/api/v1/authenticate -X POST -d "email=tu14@test14.com" -d "password=12345"
	HTTP/1.1 200 OK
	Content-Type: application/json; charset=UTF-8
	Date: Fri, 22 May 2015 16:47:47 GMT
	Content-Length: 334
	
	{
	  "account": {
	    "uid": "eccd8c58-38ec-4385-9569-6eb26a83fa17",
	    "email": "tu14@test14.com",
	    "name": "Test User 14",
	    "password": null,
	    "active": true,
	    "gravatar": null,
	    "created": "2015-05-22T15:04:45.710780544Z",
	    "updated": "2015-05-22T15:46:18.597248475Z",
	    "deleted": null
	  },
	  "status": "ok"
	}
	````

Status Codes
------------

The server returns the following status codes:

- `200`: Ok
- `201`: Created
- `400`: Bad Request
- `403`: Forbidden
- `404`: Not Found
- `500`: Internal Server Error (dont know what happened)

