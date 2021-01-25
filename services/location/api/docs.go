// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Emeric de Bernis",
            "email": "emeric.debernis@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/locations": {
            "get": {
                "description": "Get all user locations.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "locations"
                ],
                "summary": "Get locations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Location"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpapi.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "description": "Create new user locations.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "locations"
                ],
                "summary": "Create locations",
                "parameters": [
                    {
                        "description": "New location",
                        "name": "location",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateLocation"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Location"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpapi.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpapi.HTTPError"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Basic check of HTTP API health. Ensure that HTTP serving is working correctly.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "healthchecks"
                ],
                "summary": "Ping API",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "httpapi.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "HTTP status code.",
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "description": "String describing the error that occurred.",
                    "type": "string",
                    "example": "Bad Request"
                }
            }
        },
        "models.CreateLocation": {
            "type": "object",
            "required": [
                "address",
                "category_id",
                "name"
            ],
            "properties": {
                "name": {
                    "description": "Short descriptive name of the location, like \"Home\" or \"Work\".",
                    "type": "string",
                    "x-order": "1",
                    "example": "Home"
                },
                "address": {
                    "description": "Full address of the location. Should contains at least street, postal code and city.",
                    "type": "string",
                    "x-order": "2",
                    "example": "1 rue de la Poste, 75001 Paris"
                },
                "category_id": {
                    "description": "Location category foreign key.",
                    "type": "string",
                    "x-order": "3",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                }
            }
        },
        "models.Location": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "Location ID. Must be unique.",
                    "type": "string",
                    "x-order": "1",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "name": {
                    "description": "Short descriptive name of the location, like \"Home\" or \"Work\".",
                    "type": "string",
                    "x-order": "2",
                    "example": "Home"
                },
                "address": {
                    "description": "Full address of the location. Should contains at least street, postal code and city.",
                    "type": "string",
                    "x-order": "3",
                    "example": "1 rue de la Poste, 75001 Paris"
                },
                "category_id": {
                    "description": "Location category foreign key.",
                    "type": "string",
                    "x-order": "4",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "user_id": {
                    "description": "User ID. Owner of the location.",
                    "type": "string",
                    "x-order": "5",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
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
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{"http"},
	Title:       "Locations Service REST API",
	Description: "This REST API handles management of user locations. Locations can be saved in a local repository\nor fetched from third-party sources such as Google Maps \"My Places\".",
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
