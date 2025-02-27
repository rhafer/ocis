Feature: create auth-app token
  As a user
  I want to create auth-app Tokens
  So that I can use 3rd party apps

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: user creates auth-app token
    When user "Alice" creates app token with expiration time "72h" using the auth-app API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "token",
          "expiration_date",
          "created_date",
          "label"
        ],
        "properties": {
          "token": {
            "pattern": "^[a-zA-Z0-9]{16}$"
          },
          "label": {
            "const": "Generated via API"
          }
        }
      }
      """


  Scenario: user lists app tokens
    Given user "Alice" has created app token with expiration time "72h" using the auth-app API
    And user "Alice" has created app token with expiration time "72h" using the auth-app CLI
    When user "Alice" lists all created tokens using the auth-app API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "array",
        "minItems": 2,
        "maxItems": 2,
        "uniqueItems": true,
        "items": {
          "oneOf": [
            {
              "type": "object",
              "required": [
                "token",
                "expiration_date",
                "created_date",
                "label"
              ],
              "properties": {
                "token": {
                  "pattern": "^\\$2a\\$11\\$[A-Za-z0-9./]{53}$"
                },
                "label": {
                  "const": "Generated via API"
                }
              }
            },
            {
              "type": "object",
              "required": [
                "token",
                "expiration_date",
                "created_date",
                "label"
              ],
              "properties": {
                "token": {
                  "pattern": "^\\$2a\\$11\\$[A-Za-z0-9./]{53}$"
                },
                "label": {
                  "const": "Generated via CLI"
                }
              }
            }
          ]
        }
      }
      """
