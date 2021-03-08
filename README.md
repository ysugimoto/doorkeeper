# doorkeeper

Validate your pull-request with your rule and collect release notes on your deployment pull-request.

## Requirements

- Go 1.15 or later

##  Github Token

This service handles github webhook and call some github API with github token's permission and token owner's organization.

- Make sure token have `repo:status`, `repo_deployment`, and `public_repo` permission
- Make sure token owner joins organization which you want to access repositories (case of private repositories)

Github token must be set environment variable named `GITHUB_TOKEN`.

## Usage

```
go run cmd/main.go
```

Webhook server will start on "http://localhost:9000". You can change listen port by setting `PORT` environment variable.

## Webhook Secret

Regarding request comes from expected Github webhook, you need to generate secret and verify signature for each request.
Generate random secret and set environment variable named `WEBHOOK_SECRET`, and full out secret value in Github webhook secret input.

## PullRequest validation rules

Webhook server will validate your PullRequest title and description, then you can customize validation rules by putting `.doorkeeper.yml` on your repository root. Example setting file is following:

```yaml
# .doorkeeper.yml
validation:
  title:
    - kind: prefixed
      values:
        - "feat:"
        - "fix:"
  description:
    - kind: contains
      values:
        - "# Why do we need this change"
        - "# What change do you intend to"
  branches:
    - feature/*
releasenote:
  branches:
    - deployment/production
  tags:
    - v[0-9]+.[0-9]+.[0-9]+
integration:
  slack: "[slack-incoming-webhook-url]"
```

You can configure webhook behavior in root of `.doorkeeper.yml`. This file is able to have three root section

### validation field

| kind                | value type        | value description               | behaves                                                   |
|:--------------------|:------------------|:--------------------------------|:----------------------------------------------------------|
| branches            | array of string   | string or regexp string         | exact branch name or regular expression to run validation |
| title               | array of object   | have kind and value field       | validation setting for pullrequest title                  |
| description         | array of object   | have kind and value field       | validation setting for pullrequest title                  |

On `title` and `description` section, you can declare verification rule in array of `kind` and `values` object.
Validation rule configuration examples are following:


| kind      | values operator | value type                      | behaves                                                 |
|:----------|:---------------:|:--------------------------------|:--------------------------------------------------------|
| prefixed  | OR              | string                          | string MUST have prefixed words in list of values       |
| regexp    | OR              | string (valid as regexp format) | string MUST match regular expressions in list of values |
| contains  | AND             | string                          | string MUST contain all of list values                  |
| blacklist | AND             | string                          | string MUST NOT be equals to blacklist strings          |

### relasenote field

On `releasenote` field, you can declare execute making relasenote branch, tag, and integration setting.
Full validation configuration examples are following:

| kind                | value type        | value type                      | behaves                                                        |
|:--------------------|:------------------|:--------------------------------|:---------------------------------------------------------------|
| branches            | array of string   | string or regexp string         | exact brnach name or regular expression to factory releasenote |
| tags                | array of string   | regexp string                   | tag format matching regular expression                         |

### integration field

On `integration` field, you can declare notification setting for validation and releasenote process result.
Currently, only supports `slack` integration.


| integrations        | map[string]string | -                               | -                                                         |
| integrations[key]   | string            | integration type string         | Currently support `slack` only                            |
| integrations[value] | string            | integration value string        | on `slack` type, value must be valid incoming-webhook URL |

## Collection of release note

Webhook server collects release note items between base and head branch of PullRequest from pre-degined signature in order to recognize release note string.
To retrieve as note, need to contain following signature in PullRequest description:

```
<!-- RELEASE -->[Release note item here]<!-- /RELEASE -->
```

For example, description should be

```
## Why?

We need to bundle this feature for special reason.

## What?

Implement feature of this project...

## Release note

<!-- RELEASE -->
Implement special feature which user requested.
<!-- /RELEASE -->
```

Then release note is listed with `Implement special feature which user requested.` from this PullRequest.

## License

MIT

## Author

Yoshiaki Sugimoto
