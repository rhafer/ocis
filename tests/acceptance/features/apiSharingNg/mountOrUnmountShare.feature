Feature:  mount or unmount incoming share
  As a user
  I want to have control over the share received
  So that I can filter out the files and folder shared with Me

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: unmount shared resource
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Brian" unmounts share "<resource>" using the Graph API
    And user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code of responses on all endpoints should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "value"
      ],
      "properties": {
        "value": {
          "type": "array",
          "minItems": 1,
          "maxItems": 1,
          "items": {
            "type": "object",
            "required": [
              "@client.synchronize"
            ],
            "properties": {
              "@client.synchronize": {
                "const": false
              }
            }
          }
        }
      }
    }
    """
    Examples:
      | resource      |
      | textfile0.txt |
      | FolderToShare |


  Scenario Outline: mount shared resource when auto-sync is disabled
    Given user "Brian" has disabled the auto-sync share
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has created folder "folder"
    And user "Alice" has sent the following share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Brian" mounts share "<resource>" offered by "Alice" from "Personal" space using the Graph API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
        """
        {
          "type": "object",
          "required": [
            "@client.synchronize"
          ],
          "properties": {
            "@client.synchronize": {
              "const": true
            }
          }
        }
        """
    Examples:
      | resource      |
      | textfile0.txt |
      | folder        |


  Scenario Outline: enable a group share sync by only one user in a group
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Alice" has disabled the auto-sync share
    And user "Brian" has disabled the auto-sync share
    And user "Carol" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Carol" has created folder "FolderToShare"
    And user "Carol" has sent the following share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | grp1       |
      | shareType       | group      |
      | permissionsRole | Viewer     |
    When user "Alice" mounts share "<resource>" offered by "Carol" from "Personal" space using the Graph API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@client.synchronize"
        ],
        "properties": {
          "@client.synchronize": {
            "const": true
          }
        }
      }
      """
    And user "Alice" should have sync enabled for share "<resource>"
    And user "Brian" should have sync disabled for share "<resource>"
    Examples:
      | resource      |
      | textfile0.txt |
      | FolderToShare |


  Scenario Outline: disable group share sync by only one user in a group
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has created folder "FolderToShare"
    And user "Carol" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Carol" has sent the following share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | grp1       |
      | shareType       | group      |
      | permissionsRole | Viewer     |
    When user "Alice" unmounts share "<resource>" using the Graph API
    Then the HTTP status code should be "200"
    And user "Alice" should have sync disabled for share "<resource>"
    And user "Brian" should have sync enabled for share "<resource>"
    Examples:
      | resource      |
      | textfile0.txt |
      | FolderToShare |
