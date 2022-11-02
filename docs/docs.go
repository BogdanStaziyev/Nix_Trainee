// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/comments/comment/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get Comment",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments Actions"
                ],
                "summary": "Get Comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Comment"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/comments/delete/{id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete Comment",
                "tags": [
                    "Comments Actions"
                ],
                "summary": "Delete Comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/comments/save": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Save Comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments Actions"
                ],
                "summary": "Save Comment",
                "parameters": [
                    {
                        "description": "comment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.CommentRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.CommentResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/comments/update": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update Comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments Actions"
                ],
                "summary": "Update Comment",
                "parameters": [
                    {
                        "description": "comment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.CommentRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Comment"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/posts/delete/{id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete Post",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts Actions"
                ],
                "summary": "Delete Post",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/posts/post/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get Post",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts Actions"
                ],
                "summary": "Get Post",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Post"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/posts/save": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Save Post",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts Actions"
                ],
                "summary": "Save Post",
                "parameters": [
                    {
                        "description": "comment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.PostRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Post"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/v1/posts/update": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update Post",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts Actions"
                ],
                "summary": "Update Post",
                "parameters": [
                    {
                        "description": "post info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.PostRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Post"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "LoginAuth",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth Actions"
                ],
                "summary": "LoginAuth",
                "parameters": [
                    {
                        "description": "users email, users password",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.LoginAuth"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "New user registration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth Actions"
                ],
                "summary": "Register",
                "operationId": "user-register",
                "parameters": [
                    {
                        "description": "users email, users password",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RegisterAuth"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Comment": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "createdDate": {
                    "type": "string"
                },
                "deletedDate": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "postID": {
                    "type": "integer"
                },
                "updatedDate": {
                    "type": "string"
                }
            }
        },
        "domain.Post": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "createdDate": {
                    "type": "string"
                },
                "deletedDate": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "updatedDate": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "createdDate": {
                    "type": "string"
                },
                "deletedDate": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "updatedDate": {
                    "type": "string"
                }
            }
        },
        "requests.CommentRequest": {
            "type": "object",
            "required": [
                "body"
            ],
            "properties": {
                "body": {
                    "type": "string",
                    "example": "lorem ipsum"
                }
            }
        },
        "requests.LoginAuth": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@email.com"
                },
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "01234567890"
                }
            }
        },
        "requests.PostRequest": {
            "type": "object",
            "required": [
                "body",
                "title"
            ],
            "properties": {
                "body": {
                    "type": "string",
                    "example": "Lorem ipsum"
                },
                "title": {
                    "type": "string",
                    "example": "Lorem ipsum"
                }
            }
        },
        "requests.RegisterAuth": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@email.com"
                },
                "name": {
                    "type": "string",
                    "minLength": 3
                },
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "01234567890"
                }
            }
        },
        "response.CommentResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "post_id": {
                    "type": "integer"
                }
            }
        },
        "response.LoginResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "exp": {
                    "type": "integer"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "V1.echo",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "NIX TRAINEE PROGRAM Demo App",
	Description:      "REST service for NIX TRAINEE PROGRAM",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
