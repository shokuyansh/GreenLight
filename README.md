### GreenLight 
This project is an API server for movies.

Last Page Read : 235

## Features Implemented
* Json Request and Response bodies with Validation 
* Pagination,Filtering and Sorting
* Database Migrations


## Routes
| Method      | Routes           | Description  |
| ------------- |:------------- | ----- |
| GET      | /v1/healthcheck | End point for checking health status of server |
| GET      | /v1/movies   |   End point for viewing a list of movies |
| GET      | /v1/movies/:id     |   End point for viewing a particular movie |
|PATCH       | /v1/movies/:id     | End point for updating a particular movie |
|DELETE       | /v1/movies/:id     | End point for deleting a particular movie |
| POST | /v1/movies      | End point for creating a movie    |

