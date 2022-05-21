package verify

func FilterSchema() string {
	return `{
	"$schema": "http://json-schema.org/draft-07/schema",
	"title": "Aidoku filters specification",
	"definitions": {
		"arrayOfStrings": {
			"type": "array",
			"items": {
				"type": "string"
			}
		},
		"titleOrAuthorFilter": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": [
						"title",
						"author"
					]
				}
			},
			"additionalProperties": false
		},
		"selectFilter": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": ["select"]
				},
				"name": {
					"type": "string"
				},
				"options": {
					"$ref": "#/definitions/arrayOfStrings"
				},
				"default": {
					"type": "integer"
				}
			},
			"additionalProperties": false
		},
		"checkFilter": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": ["check"]
				},
				"name": {
					"type": "string"
				},
				"default": {
					"type": "boolean"
				},
				"canExclude": {
					"type": "boolean"
				}
			},
			"additionalProperties": false
		},
		"genreFilter": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": ["genre"]
				}, 
				"name": {
					"type": "string"
				},
				"canExclude": {
					"type": "boolean"
				},
				"id": {
					"type": "string"
				},
				"default": {
					"type": "boolean"
				}
			},
			"additionalProperties": false
		},
		"sortFilter": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": ["sort"]
				},
				"name": {
					"type": "string"
				},
				"canAscend": {
					"type": "boolean"
				},
				"options": {
					"$ref": "#/definitions/arrayOfStrings"
				},
				"default": {
					"type": "object",
					"properties": {
						"index": {
							"type": "integer"
						},
						"ascending": {
							"type": "boolean"
						}
					},
					"additionalProperties": false
				}
			},
			"additionalProperties": false
		},
		"groupFilter": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": ["group"]
				},
				"name": {
					"type": "string"
				},
				"filters": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/filter"
					}
				}
			},
			"additionalProperties": false
		},
		"filter": {
			"anyOf": [
				{
					"$ref": "#/definitions/checkFilter"
				},
				{
					"$ref": "#/definitions/genreFilter"
				},
				{
					"$ref": "#/definitions/titleOrAuthorFilter"
				},
				{
					"$ref": "#/definitions/selectFilter"
				},
				{
					"$ref": "#/definitions/groupFilter"
				},
				{
					"$ref": "#/definitions/sortFilter"
				}
			]
		}
	},
	"type": "array",
	"items": {
		"$ref": "#/definitions/filter"
	}
}
`
}

func SettingsSchema() string {
	return `{
	"$schema": "http://json-schema.org/draft-07/schema",
	"title": "Aidoku settings specification",
	"definitions": {
		"any": {
			"anyOf": [
				{
					"type": "array",
					"items": {
						"$ref": "#/definitions/any"
					}
				},
				{
					"type": "boolean"
				},
				{
					"type": "integer"
				},
				{
					"type": "null"
				},
				{
					"type": "object",
					"additionalProperties": {
						"$ref": "#/definitions/any"
					}
				},
				{
					"type": "string"
				}
			]
		},
		"settings": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"enum": [
						"group",
						"select",
						"multi-select",
						"switch",
						"stepper",
						"segment",
						"text",
						"page",
						"button",
						"link"
					]
				},
				"key": {
					"type": "string"
				},
				"action": {
					"type": "string"
				},
				"title": {
					"type": "string"
				},
				"subtitle": {
					"type": "string"
				},
				"footer": {
					"type": "string"
				},
				"placeholder": {
					"type": "string"
				},
				"values": {
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"titles": {
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"default": {
					"$ref": "#/definitions/any"
				},
				"notification": {
					"type": "string"
				},
				"requires": {
					"type": "string"
				},
				"requiresFalse": {
					"type": "string"
				},
				"minimumValue": {
					"type": "number"
				},
				"maximumValue": {
					"type": "number"
				},
				"destructive": {
					"type": "boolean"
				},
				"external": {
					"type": "boolean"
				},
				"items": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/settings"
					}
				},
				"autocapitalizationType": {
					"type": "integer"
				},
				"autocorrectionType": {
					"type": "integer"
				},
				"spellCheckingType": {
					"type": "integer"
				},
				"keyboardType": {
					"type": "integer"
				},
				"returnKeyType": {
					"type": "integer"
				}
			},
			"required": [
				"type"
			]
		}
	},
	"type": "array",
	"items": {
		"$ref": "#/definitions/settings"
	}
}
`
}

func SourceSchema() string {
	return `{
	"$schema": "http://json-schema.org/draft-07/schema",
	"title": "Aidoku source specification",
	"definitions": {
		"language": {
			"type": "object",
			"properties": {
				"code": {
					"type": "string"
				},
				"value": {
					"type": "string"
				},
				"default": {
					"type": "boolean"
				}
			},
			"required": [
				"code"
			],
			"additionalProperties": false
		},
		"sourceInfo": {
			"type": "object",
			"properties": {
				"id": {
					"type": "string"
				},
				"lang": {
					"type": "string"
				},
				"name": {
					"type": "string"
				},
				"version": {
					"type": "integer"
				},
				"url": {
					"type": "string"
				},
				"urls": {
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"nsfw": {
					"type": "integer"
				}
			},
			"required": [
				"id",
				"lang",
				"name",
				"version"
			],
			"additionalProperties": false
		},
		"listing": {
			"type": "object",
			"properties": {
				"name": {
					"type": "string"
				},
				"flags": {
					"type": "integer"
				}
			},
			"required": [
				"name"
			],
			"additionalProperties": false
		} 
	},
	"type": "object",
	"properties": {
		"info": {
			"$ref": "#/definitions/sourceInfo"
		},
		"languages": {
			"type": "array",
			"items": {
				"$ref": "#/definitions/language"
			}
		},
		"listings": {
			"type": "array",
			"items": {
				"$ref": "#/definitions/listing"
			}
		}
	},
	"required": [
		"info"
	]
}
`
}
