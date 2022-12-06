# Overview
It's an API on Echo framework.
My aim is to work with this framework and personal development.
There is a useful set of tools, which are described below.
(I know that I can't store environment variables in the public domain, I need to put them in github secrets, but otherwise you won't be able to fully test my application)

## Inside:

- Registration and Authentication
- Registration and Authentication with Google
- CRUD API for posts, comments
- Migrations
- Request validation
- Swagger docs
- Environment configuration
- Docker development environment

## Test Routes
- I uploaded the finished api to the docker hub. By running the docker compose file in the root directory you can test it.
- I have also attached a postman file for your convenience.
- There is also a swagger specification you can run it in: http://localhost:8080/swagger/

## Endpoints


- CHECK PING GET http://localhost:8080/api/v1


- Registration POST http://localhost:8080/register
- Authentication POST http://localhost:8080/login
- Google Authentication GET http://localhost:8080/auth/google/login


- Swagger GET http://localhost:8080/swagger/


- Save POSTS POST http://localhost:8080/api/v1/posts/save
- Get POSTS GET http://localhost:8080/api/v1/posts/post/{id}
- Update POSTS http://localhost:8080/api/v1/posts/update/{id}
- Delete POSTS http://localhost:8080/api/v1/posts/delete/{id}


- Save COMMENTS POST http://localhost:8080/api/v1/comments/save/{postID}
- Get COMMENTS GET http://localhost:8080/api/v1/comments/comment/{id}
- Update COMMENTS http://localhost:8080/api/v1/comments/update/{id}
- Delete COMMENTS http://localhost:8080/api/v1/comments/delete/{id}
