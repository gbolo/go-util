// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "dev@appname"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/client": {
            "get": {
                "description": "Returns a list of all clients",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clients"
                ],
                "summary": "Returns a list of all clients",
                "responses": {
                    "200": {
                        "description": "a list of clients",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/backend.Client"
                            }
                        }
                    },
                    "500": {
                        "description": "an error occurred. Usually a database issue",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an existing Client with the specified ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clients"
                ],
                "summary": "Updates an existing Client",
                "parameters": [
                    {
                        "description": "Add Client",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/backend.Client"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "client has been updated",
                        "schema": {
                            "$ref": "#/definitions/backend.successResponse"
                        }
                    },
                    "304": {
                        "description": "client does not need updating, it's already at that state",
                        "schema": {
                            "$ref": "#/definitions/backend.successResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request body",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    },
                    "404": {
                        "description": "a client with that ID does not exist",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    },
                    "500": {
                        "description": "server was unable to process the request. Usually a database issue",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new Client",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clients"
                ],
                "summary": "Add a new Client",
                "parameters": [
                    {
                        "description": "Add Client",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/backend.Client"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "new client was added",
                        "schema": {
                            "$ref": "#/definitions/backend.successResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request body",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    },
                    "409": {
                        "description": "a client with that ID already exists. Try running an update instead",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    },
                    "500": {
                        "description": "server was unable to process the request. Usually a database issue",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    }
                }
            }
        },
        "/v1/client/{id}": {
            "get": {
                "description": "Return the status of a client with specified ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clients"
                ],
                "summary": "Return the status of a client",
                "responses": {
                    "200": {
                        "description": "returns a client's status",
                        "schema": {
                            "$ref": "#/definitions/backend.ClientStatus"
                        }
                    },
                    "404": {
                        "description": "a client with that ID does not exist",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    },
                    "500": {
                        "description": "an error occurred. Usually a database issue",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a Client by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clients"
                ],
                "summary": "Delete a Client by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Client ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "client was deleted or does not exist",
                        "schema": {
                            "$ref": "#/definitions/backend.successResponse"
                        }
                    },
                    "500": {
                        "description": "an error occurred. Usually a database issue",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    }
                }
            }
        },
        "/v1/healthz": {
            "get": {
                "description": "If the status is not 200, then the application is unhealthy",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Returns the health of the application",
                "responses": {
                    "200": {
                        "description": "server is healthy",
                        "schema": {
                            "$ref": "#/definitions/backend.successResponse"
                        }
                    },
                    "500": {
                        "description": "server is unhealthy. Usually a database issue",
                        "schema": {
                            "$ref": "#/definitions/backend.errorResponse"
                        }
                    }
                }
            }
        },
        "/v1/version": {
            "get": {
                "description": "Returns version information",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Returns version information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/backend.versionInfo"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "backend.Client": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "backend.ClientStatus": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "reachable": {
                    "type": "boolean"
                },
                "status": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "backend.errorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "backend.successResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "backend.versionInfo": {
            "type": "object",
            "properties": {
                "build_date": {
                    "type": "string"
                },
                "build_ref": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "appname",
	Description: "Swagger API appname",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
