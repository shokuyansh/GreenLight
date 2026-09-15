### GreenLight 
This project is an API server for movies.

Last Page Read : 484

## Features Implemented
* Json Request and Response bodies with Validation 
* Pagination,Filtering and Sorting
* Database Migrations
* Rate limiting
* MiddleWares including (panicRecovery, enableCORS, authentication,... etc)
* GraceFul Shutdown
* User Activation via activation tokens
* Authentication via stateful auth tokens
* Authorization via User Permissions 
* Enabling Cross-Origin Requests 
* Automatic Versioning
* Quality Control
* Self-Contained Build



## Routes
| Method      | Routes           | Description  |
| ------------- |:------------- | ----- |
| GET | /debug/vars/ | Display Application Metrics
| GET      | /v1/healthcheck | Check health status of server |
| GET      | /v1/movies   |   View a list of movies |
| GET      | /v1/movies/:id     |   View a particular movie |
|PATCH       | /v1/movies/:id     | Updating a particular movie |
|DELETE       | /v1/movies/:id     | Delete a particular movie |
| POST | /v1/movies      | Create a new movie    |
| POST     | /v1/users     |   Register a new user(read perm granted & email is sent in background) |
| PUT   | /v1/users/activated | Activate a user account with activation token |
| POST | /v1/tokens/activation | Re-send a new activation token to user via email |
|POST | /v1/tokens/authentication | Create authentication token for an existing user |


