# caddyslack [![Build Status](https://travis-ci.org/zahlz/SwiftPasscodeLock.svg?branch=master)](https://travis-ci.org/zahlz/SwiftPasscodeLock) [![Go Report Card](https://goreportcard.com/badge/github.com/zahlz/caddyslack)](https://goreportcard.com/report/github.com/zahlz/caddyslack)

Caddy plugin to filter and relay incoming WebHook requests to slack

```
slack [endpoint] {
  url https://hooks.slack.com/services/ID/ID/TOKEN
  [only]
    [json.field.to.keep]
    [json.field.to.keep.as.well]
  [delete]
    [json.field.to.delete]
    [json.field.to.delete.as.well]
}
```

![](/doc/caddySlack.png)

## Examples
### Full Example with ratelimit
[ratelimit caddy plugin](https://caddyserver.com/docs/ratelimit)
```
slack /toSlack {
  url https://hooks.slack.com/services/ID/ID/TOKEN
  only
    text
}
ratelimit /toSlack 2 3
```
The endpoint `/toSlack` accepts 2 or 3 (with burst) `POST` requests per second. The json-body will be filtered for everything except `text`, and forwarded to `https://hooks.slack.com/services/ID/ID/TOKEN`

### delete

Caddyfile
```
slack /toSlack {
  url https://hooks.slack.com/services/ID/ID/TOKEN
  delete
    attachments.title
}
```

A `POST` request to `/toSlack` with the following body

```json
{
  "text": "Hello",
  "attachments": [
        {
            "title": "App hangs on reboot",
            "text": "If I restart my computer without quitting your app, it stops the reboot sequence.",
        }
    ]
}
```

will be forwarded to `https://hooks.slack.com/services/ID/ID/TOKEN` as

```json
{
  "text": "Hello",
  "attachments": [
        {
            "text": "If I restart my computer without quitting your app, it stops the reboot sequence.",
        }
    ]
}
```


### only

Caddyfile
```
slack /toSlack {
  url https://hooks.slack.com/services/ID/ID/TOKEN
  only
    text
    icon
}
```

A `POST` request to `/toSlack` with the following body

```json
{
  "text": "Hello",
  "channel": "notallowed",
  "icon": "ghost"
}
```

will be forwarded to `https://hooks.slack.com/services/ID/ID/TOKEN` as

```json
{
  "text": "Hello",
  "icon": "ghost"
}
```
