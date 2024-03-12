// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Maintainer",
            "email": "v.rianov@kabanov.agency"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/places": {
            "get": {
                "description": "Get places",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "places"
                ],
                "summary": "Get places",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Latitude",
                        "name": "latitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude",
                        "name": "longitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Latitude delta",
                        "name": "latitudeDelta",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude delta",
                        "name": "longitudeDelta",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Labels",
                        "name": "labels",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Place"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "put": {
                "description": "Update place",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "places"
                ],
                "summary": "Update place",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Place",
                        "name": "place",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdatePlaceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Place updated"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "post": {
                "description": "Create place",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "places"
                ],
                "summary": "Create place",
                "parameters": [
                    {
                        "description": "Place",
                        "name": "place",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreatePlaceRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Place created"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/places/{id}": {
            "get": {
                "description": "Get place",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "places"
                ],
                "summary": "Get place",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Place ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Place"
                        }
                    },
                    "404": {
                        "description": "Place not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "delete": {
                "description": "Delete place",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "places"
                ],
                "summary": "Delete place",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Place ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Place deleted"
                    },
                    "404": {
                        "description": "Place not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Coordinates": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "models.CreatePlaceRequest": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "closed_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "labels": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "open_at": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "models.Place": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "coordinates": {
                    "$ref": "#/definitions/models.Coordinates"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "instagram": {
                    "type": "string"
                },
                "labels": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "locationURL": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "models.UpdatePlaceRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "labels": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "34.159.168.142",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Favs API",
	Description:      "This is a documentation for favs API endpoints.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
