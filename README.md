# Foody
I know, the name of the repository is not the same as the name of the API, but in the progress of creation I came up with some ideas and thought “This name is going to be good”.
So Foody was just born.

### Instalation
Run the make file with the following command to install all the package, run the test's and check the files:
```CMD
make audit
```
### How to Run it:
Type in youre terminal before clone the repository, at this moment I asume you have git in youre machine so:
```CMD
make run
```
#### Important
Foody is using koyeb a service who provides an alternative to use the DB in your local machine if you dont want to use it you need to create a DB in youre local machine,remember do the next sql migrations:
```CMD
migrate create -seq -ext=.sql -dir=./migrations create_foods_table
```
This returns the creation of two files:
```
./migrations/
├── 000001_create_foods_table.down.sql
└── 000001_create_foods_table.up.sql
```
into the file 000001_create_foods_table.down.sql, copy the follow query:
```SQL
CREATE TABLE IF NOT EXISTS foods (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
title text NOT NULL,
types text[] NOT NULL,
version integer NOT NULL DEFAULT 1
);
```
and in the 000001_create_foods_table.up.sql the next:
```SQL
DROP TABLE IF EXISTS foods;
```
now is time to run the migration (this is a brew example to how to do the migration to remote services):
```CMD
migrate -source="s3://<bucket>/<path>" -database=$EXAMPLE_DSN up
migrate -source="github://owner/repo/path#ref" -database=$EXAMPLE_DSN up
migrate -source="github://user:personal-access-token@owner/repo/path#ref" -database=$EXAMPLE_DSN up
```
So for now you got the foods table if you following the steps, now run the index:
```CMD
migrate create -seq -ext .sql -dir ./migrations add_food_indexes
```
in the file 000002_add_foods_indexes.up.sql copy the following code:
```SQL
CREATE INDEX IF NOT EXISTS foods_title_idx ON foods USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS fooods_types_idx ON foods USING GIN (types);
```
same thing with the 000002_add_foods_indexes.down.sql:
```SQL
DROP INDEX IF EXISTS foods_title_idx;
DROP INDEX IF EXISTS foods_types_idx;
```
Now execute the up migration to add the indexes to DB

### To create the tables for the user auth:
At this point youre migrate two migration files to youre DB, so for now I asume you now how to create and send the migratios so I only focus in what the file have
#### create_users_table_up.sql
```SQL
CREATE TABLE IF NOT EXISTS users (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
name text NOT NULL,
email citext UNIQUE NOT NULL,
password_hash bytea NOT NULL,
activated bool NOT NULL,
version integer NOT NULL DEFAULT 1
);
```
#### create_user_table_down.sql
```SQL
DROP TABLE IF EXISTS users;
```
#### create_tokens_table_up.sql
```SQL
CREATE TABLE IF NOT EXISTS tokens (
hash bytea PRIMARY KEY,
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
expiry timestamp(0) with time zone NOT NULL,
scope text NOT NULL
);
```
#### create_tokens_table.down.sql
```SQL
DROP TABLE IF EXISTS tokens;
```
#### add_permissions.up.sql
```SQL
CREATE TABLE IF NOT EXISTS permissions (
id bigserial PRIMARY KEY,
code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
PRIMARY KEY (user_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
('foods:read'),
('foods:write');
```
#### add_permissions.down.sql
```SQL
DROP TABLE IF EXISTS users_permissions;
DROP TABLE IF EXISTS permissions;
```
Now at this point the API is setting up to be use it
--- 
# Foody API 
Let's begin talking about the end point's of the API. Now it have 10 enpoints which are:

## This Enpoinst are the create, delete, update, and read of table "foods" |The endpoint named "healthcheck" is an exception|

| Method  | EndPoint | Description |
|-------- |----------|-------------|
|  GET |  "/v1/foods" |  healthcheck: This endpoint returns the live status of the API, version, date, etc.| 
|  POST | "/v1/foods"  | createFood: This endpoint create an record into the DB.|
|  GET | "/v1/foods"  | listFood: This enpoint list all the records on the DB.|
|  GET | "/v1/foods/:id"  | showFood: This enpoint returns one record, searching the record by the ID.|
|  PATCH  | "/v1/foods/:id"  | updateFood: This endpoint updates a record using the ID as the showFood.|
|  DELETE  | "/v1/foods/:id"  | deleteFood: This endpoint delete a record using the ID.|

## The use of the next enpoints are the registration, autentication and activation of user's
| Method  | EndPoint | Description |
|-------- |----------|-------------|
|  POST  | "/v1/users"  | registerUser: This enpoint register one user into the DB.|
|  PUT  | "/v1/users/activated"  | activateUser: This enpoints use a token generated by the registration to authenticate the user|
|  POST    | "/v1/tokens/authentication"  | createAuthentication: This enpoints generate a token wich has the utility for use the enpoints mentioned befor (create, delete, update in the table food)|

## About the last one is the stats of the API 
| Method  | EndPoint | Description |
|-------- |----------|-------------|
|  GET  | "/v1/debug/vars"  | expVar: This enpoint use a expvar go package to show the stats of the API|

---
Well, about the use of the API I take into account that you are thinking “well, how do I use it?” so in this part I describe how to make the request to use it

# AUTH

#### Register and user into the API
Well to describe the process behind of this request, when the user is register into the API is launching a goroutine wich send a SMTP using parameter email in the JSON
```CMD
    BODY='{"name": "Tenki Me", "email": "tenkivale@example.com", "password": "pa55word"}'
    curl -i -d "$BODY" localhost:4000/v1/users
```
#### Activate user
To activate a user it is important to complete the above steps before continuing. When a user registers in the database, a token is assigned to him/her that will be sent to him/her by mail once the registration process is finished. I forgot to mention that when a user 
is registered and activated, the API gives him/her a “role” that by default is read -> I will go deeper into this topic below
```CMD
  curl -X PUT -d '{"token": "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}' localhost:4000/v1/users/activated
```
#### Authenticate user
After the previous steps ther user is register and activated in to the API, so for now use the credentials will passed before in the registration process to authenticate the user.
```CMD
  BODY='{"email": "tenkivale@example.com", "password": "pa55word"}'
  curl -d "$BODY" localhost:4000/v1/tokens/authentication
```
-> API response:
```JSON
  {
  	"authentication_token": {
  		"token": "WRLR2UTQ6X35XVK6PVR3KG2ABA", -> This is an example of the token
  		"expiry": "2025-02-14T13:36:57.732112981-05:00" -> The date when the token expire
  	}
  }
```
---
## At this moment you can use the rest of the endpoints if youre following the previous steps
Do you remember now the token generated in the authentication step? It is the most important thing because without it you cannot use the following endpoints
#### CreateFood:
This enpoint has a specially thing beacuse only user with the permission "read/write" can create new records into the foods table
```CMD
  BODY='{"title":"Soup", "types":["onion", "potatoo", "pepper"]}'
  curl -i -d "$BODY" -H "Authorization: Bearer WRLR2UTQ6X35XVK6PVR3KG2ABA" localhost:4000/v1/foods
```
If you dont have the permissions to create records in to the table, the API returns the error below:
```JSON
HTTP/1.1 403 Forbidden
Access-Control-Allow-Origin: *
Content-Type: application-json
Vary: Authorization
Date: Thu, 13 Feb 2025 18:57:29 GMT
Content-Length: 97

{
	"error": "your user account doesn't have the necessary permissions to access this resource"
}
```
instead return this:
```JSON
HTTP/1.1 201 Created
Access-Control-Allow-Origin: *
Content-Type: application-json
Location: /v/foods/5
Vary: Authorization
Date: Thu, 13 Feb 2025 19:07:53 GMT
Content-Length: 199

{
	"food": {
		"ID": 5,
		"CreateAt": "2025-02-13T19:07:54Z",
		"Title": "Soup",
		"Types": [
			"onion",
			"potatoo",
			"pepper"
		],
		"Version": 1
	}
}
{Title:Soup Types:[onion potatoo pepper]}
```
#### updateFood
Likewise I told before only user with the permissions of write and read can update records.
```CMD
  BODY='{"title":"Soup Update", "types":["onion", "update", "pepper"]}'
  curl -X PATCH -d "$BODY" -H "Authorization: Bearer ITVDVHAJEJCXTKKDRDUZ3F5IXU" localhost:4000/v1/foods/5
```
if everthing is good, the API returns this:
```JSON
{
	"food": {
		"ID": 5,
		"CreateAt": "2025-02-13T19:07:54Z",
		"Title": "Soup Update", -> The update is present
		"Types": [
			"onion",
			"update", -> And here
			"pepper"
		],
		"Version": 2 -> The version number is higher because the registration has differences
	}
}
```
#### showFood
This endpoint can access by all the users (obviosly register and actived) just because only need the "read" permission to use it.
```CMD
  curl -H "Authorization: Bearer ITVDVHAJEJCXTKKDRDUZ3F5IXU" localhost:4000/v1/foods/5
```
API returns:
```JSON
{
	"food": {
		"ID": 5,
		"CreateAt": "2025-02-13T19:07:54Z",
		"Title": "Soup Update",
		"Types": [
			"onion",
			"update",
			"pepper"
		],
		"Version": 2
	}
}
```
#### listFood
This endpoint like I told in the previous step, is acceded by all the user register and activated in the API.
```CMD
  curl -H "Authorization: Bearer ITVDVHAJEJCXTKKDRDUZ3F5IXU" localhost:4000/v1/foods
```
The API returns:
```JSON
{
	"foods": [
		{
			"ID": 3,
			"CreateAt": "2025-02-05T15:05:17Z",
			"Title": "Pizza",
			"Types": [
				"pepperoni",
				"fastfood",
				"chesee"
			],
			"Version": 1
		},
		{
			"ID": 4,
			"CreateAt": "2025-02-05T15:05:52Z",
			"Title": "Soup Test",
			"Types": [
				"onion",
				"potatoo",
				"pepper"
			],
			"Version": 1
		},
		{
			"ID": 5,
			"CreateAt": "2025-02-13T19:07:54Z",
			"Title": "Soup Update",
			"Types": [
				"onion",
				"update",
				"pepper"
			],
			"Version": 2
		}
	],
	"metada": {
		"CurrentPage": 1,
		"PageSize": 20,
		"FirstPage": 1,
		"LastPage": 1,
		"TotalRecords": 3
	}
}
```
---
For now feel safe to clone this repository and run in youre local machine, is not deployed now but in a few days I guess it happens
