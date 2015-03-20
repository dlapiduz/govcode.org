# govcode API Reference
Government Open Source Project Explorer
v2.0.1

## Resources
###/repos/
Github repositories related to goverment open source projects.

#### Parameters
| Name    | Located in | Description                           | Required | Type   |
|---------|------------|---------------------------------------|----------|--------|
| perPage | query      | Amount of results to return per page. | No       | number |

#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | An array of repositories | [[repository](#schema_repository)] |
| default | json | Unexpected error | [error](#schema_error) |

###/repos/:name
One Github repository with the name of the `:name` parameter.

#### Parameters
None


#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | A single repository | [repository](#schema_repository) |
| default | json | Unexpected error | [error](#schema_error) |


###/orgs/
Github organizations writing code for goverment open source projects.

#### Parameters
None

#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | An array of repositories | [[organization](#schema_organization)] |
| default | json | Unexpected error | [error](#schema_error) |

###/orgs/:name
Github organization whose name matches the `:name` parameter.

#### Parameters
None

#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | An array of repositories | [organization](#schema_organization) |
| default | json | Unexpected error | [error](#schema_error) |


###/users/
Github users writing code for goverment open source projects.

#### Parameters
None

#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | An array of repositories | [[user](#schema_user)] |
| default | json | Unexpected error | [error](#schema_error) |

###/users/:id
Github user whose id matches the `:id` parameter.

#### Parameters
None

#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | An array of repositories | [user](#schema_user) |
| default | json | Unexpected error | [error](#schema_error) |


###/stats/


###/issues
Github issues on projects related to government open source projects.

#### Parameters
| Name    | Located in | Description                           | Required | Type   |
|---------|------------|---------------------------------------|----------|--------|
| perPage | query      | Amount of results to return per page. | No       | number |
| page | query      | The page number to query | No       | number |
| languages | query      | The coding language for the project | No       | string |
| organizations | query      | The organization that owns the issue | No       | string |
| state | query      | The issue state, either `open` or `closed` | No       | string |
| label | query      | One of the labels for the issue | No       | string |

#### Responses
| Code    | Type | Description              | Schema            |
|---------|------|--------------------------|-------------------|
| 200 | json | An array of repositories | [[issue](#schema_issue)] |
| default | json | Unexpected error | [error](#schema_error) |

#### Example
```
/issues/?perPage=10&page=2&repoId=20&orgId=200&languages=python&organizations=18F&state=open&label=help
```

## Schemas
<a name="schema_repository"></a>
### repository
```json
{
  "properties": {
	  "Id": {
	    description: "A unique identifier representing the specific repository"
	    type: "number",
	    format: "int"
	  },
	  "GhId": {
	    description: "The Github-specific ID for the repository",
	    type: "number"
	    format: "int"
	  },
	  "Name": {
	    description: "The name assigned to the repository in Github",
	    type: "string"
	  },
	  "Forks": {
	    description: "The number of repository 'forks' on Github",
	    type: "number",
	    format: "int"
	  },
	  "Watchers": {
	    description: "The number of people watching the repository on Github",
	    type: "number",
	    format: "int"
	  },
	  "Stargazers": {
	    description: "The number of people who've starred the respository on Github",
	    type: "number",
	    format: "int"
	  },
	  "Size": {
	    description: "I don't know",
	    type: "number",
	    format: "int"
	  },
	  "OpenIssues": {
	    description: "The number of open issues",
	    type: "number",
	    format: "int"
	  },
	  "Description": {
	    description: "Description of the repository",
	    type: "string"
	  },
	  "Language": {
	    description: "The programming language the repository is written in",
	    type: "string",
	    example: "c++"
	  },
	  "LastCommit": {
	    description: "Time of the last time code was committed to the repository",
	    type: ([**schema/timestamp**](#schema_timestamp))
	  },
	  "LastPull": {
	    description: "Time of the pull request being merged into the repository",
	    type: ([**schema/timestamp**](#schema_timestamp))
	  },
	  "CommitCount": {
	    description: "The number of total commits for the repository",
	    type: "number",
	    format: "int"
	  },
	  "OrganizationId": {
	    description: "The ID of the organization that the repository belongs to on Github",
	    type: "string"
	  },
	  "OrganizationLogin": {
	    description: "The login name of the organization that the repository belongs to on Github",
	    type: "string"
	  },
	  "DaysSincePull": {
	    description: "Number of days since the last pull request was closed on the repository",
	    type: "number",
	    format: "int"
	  },
	  "DaysSinceCommit": {
	    description: "Number of days since the last commit on the repository"
	    type: "number"
	    format: "int"
	  },
	  "Commits": {
	    description: "The commits in the repo not currently in use",
	    default: null
	  },
	  "Pulls": {
	    description: "The pull requests in the repo, not currently in use",
	    default: null
	  },
	  "Organization": {
	    description: "The organization that owns the repo",
	    type: ([**schema/organization**](#schema_organization))

	  "RepoStat": {
	    description: "Not in use",
	    default: null
	  },
	  "Ignore": {
	    description: "Whether the current repository is being ignored by the Github user",
	    type: "boolean"
	  },
	  "GhCreatedAt":{  
	    description: "When the repository was created in Github",
	    type: ([**schema/timestamp**](#schema_timestamp))
	  },
	  "GhUpdatedAt":{  
	  	 description: "When the repository was updated in Github",
	    type: ([**schema/timestamp**](#schema_timestamp))
	  },
	  "CreatedAt": {
	    description: "The datetime when the repository was created",
	    type: "string"
	    format: "utc"
	  },
	  "UpdatedAt": {
	  	 description: "The datetime when the repository was updated",
	    type: "string"
	    format: "utc"
	  },
	  "HelpWantedIssueCount": {
	    description: "Number of issues with a 'help wanted' tag on them"
	    type: "number"
	    format: "int"
	  }
  }
}
```

<a name="schema_organization"></a>
### Organization
A Github organization.

```json
{
  "properties": {
	  "Id": {
	    description: "A unique identifier representing the specific organization"
	    type: "number",
	    format: "int"
	  },
	  "Name": {
	    description: "The name of the organization",
	    type: "string"
	  },
	  "Login": {
	    description: "The login of the organization, not in use",
	    default: ""
	  },
	  "Ignore": {
	    description: "Whether the organization is being ignored on Github or not",
	    type: "boolean"
	  },
	  "Repositories": {
	    description: "A list of repositories owned by the organization",
	    default: null
	  },
	  "CreatedAt": {
	    description: "The datetime when the organization was created",
	    type: "string"
	    format: "utc"
	  }
	}
}
```

<a name="schema_user"></a>
### User
A Github user

```json
{
  "properties": {
    "Id": {
      description: "A unique identifier for the Github user",
      type: "number"
      format: "int"
    },
    "GhId": {
      description: "A unique identifier for the Github user from the Github site",
      type: "number",
      format: "int"
    },
    "Login": {
      description: "The login name of the user",
      type: "string"
    },
    "AvatarUrl": {
      description: "A url of the avatar image for the user",
      type: "string"
    },
    "CommitCount": {
      description: "The amount of commits for the user",
      type: "number",
      format: "int"
    },
    "OrgList": {
      description: "A string list of organizations the user belongs to",
      example: "{usnationalarchives,18f}",
      type: "string"
    },
    "Commits": {
      description: "Commits committed by the user, not in use",
      default: null
    }
  }
}
```


<a name="schema_issue"></a>
### Issue
A Github issue

```json
{
  properties: {  
    "Id": {
      description: "A unique identifier for the Github user",
      type: "number"
      format: "int"  
    },
    "RepositoryId": {
      description: "The unique identifier of the repository the issue belongs to",
      type: "number",
      format: "int"
    },
    "Number": {
    
    },
    "Title": {
      description: "The title of the issue",
      type: "string"
    },
    "Body": {
      description: "The body text of the issue",
      type: "string",
      format: "Github-flavored markdown"
    },
    "Url": {
      description: "The url of the issue",
      type: "string"
      example: "https://github.com/18F/C2/issues/247"
    },
    "Labels": {
      description: "The labels applied to the issue",
      type: "string",
      example: "help wanted, frontend"
    },
    "State": {
      description: "The state the issue is currently in",
      example: "closed",
      type: "string"
    },
    "OrganizationLogin": {
      description: "The organization's login name that owns the repository that the issue belongs to",
      type: "string"
    },
    "Language": {
      description: "The language that the repository is mainly composed of that the issue belongs to",
      example: "ruby",
      type: "string"
    },
    "RepositoryName": {
      description: "The name of the repository which the issue belongs to",
      type: "string"
    },
    "GhCreatedAt": {  
	   description: "When the issue was created in Github",
	   type: ([**schema/timestamp**](#schema_timestamp))
    },
    "GhUpdatedAt": {  
	   description: "When the issue was last updated in Github",
	   type: ([**schema/timestamp**](#schema_timestamp))
    },
    "GhClosedAt": {  
	   description: "When the issue was closed in Github",
	   type: ([**schema/timestamp**](#schema_timestamp))
    },
    "CreatedAt": {
      description: "When the issue was created",
      type: "string"
    }
  }
}
```


<a name="schema_timestamp"></a>
### Timestamp
A string representing a time and date.

```json
{
  properties: {
    "Time": {
      description: "UTC representation of a time and date",
      example: "0001-01-01T00:00:00Z",
      type: "string"
    },
    "Valid": {
      description: "Whether the date is a real date or not. If not valid, do not use the date",
      type: "boolean"
    }
  }
}
```

<a name="schema_error"></a>
### Error
```json
{
  properties: {
  }
}
```