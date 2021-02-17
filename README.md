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

## PullRequest validation rules

Webhook server will validate your PullRequest title and description, then you can customize validation rules by putting `.pullrequest.yml` on your repository root. Example setting file is following:

```yaml
# .pullrequest.yml
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
```

`.pullrequest.yml` can have a couple of root section -- `title` and `description` -- these correspond to PullRequest title and description.
You can declare verification rule in array of `kind` and `values` object. We show all enable configurations following:

| kind      | values operator | value type                      | behaves                                                 |
|:----------|:---------------:|:--------------------------------|:--------------------------------------------------------|
| prefixed  | OR              | string                          | string MUST have prefixed words in list of values       |
| regexp    | OR              | string (valid as regexp format) | string MUST match regular expressions in list of values |
| contains  | AND             | string                          | string MUST contain all of list values                  |
| blacklist | AND             | string                          | string MUST NOT be equals to blacklist strings          |


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

And this service can determines branch of collect relase note branch, it means `your release branch ` of webhook path following `/webhook/`.
If you set webhook URL as `https://[TBD]/webhook/deployment/production`, release note will collect on `deployment/production` branch.

You may want to release multiple apllication from single repository (e.g. monorepo), webhook allows to multiple branches by setting path as regex style.
For example, If you set webhook URL as `https://[TBD]webhook/deployment/*`, we compares target branch with regular expression so you can make release note for `deployment/service-1`, `deployment/service-2` branches PullRequest.
