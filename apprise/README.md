# Apprise

A Dagger module for [Apprise](https://github.com/caronc/apprise) - send notifications to all of the most popular notification services.

## Usage

### Send a notification

```bash
dagger call --mod ./apprise send \
  --title "Test Title" \
  --body "Hello from Apprise" \
  --service env:APPRISE_SERVICE_URL \
  stdout
```

### Using with environment variables

```bash
export APPRISE_SERVICE_URL="syslog://"

dagger call --mod ./apprise send \
  --title "Build Complete" \
  --body "Deployment successful" \
  --service env:APPRISE_SERVICE_URL \
  stdout
```

### Using with Discord webhook

```bash
export DISCORD_WEBHOOK="discord://WEBHOOK_ID/WEBHOOK_TOKEN"

dagger call --mod ./apprise send \
  --title "Alert" \
  --body "Something happened" \
  --service env:DISCORD_WEBHOOK \
  stdout
```

### Supported services

Apprise supports many notification services including:

- Discord: `discord://WEBHOOK_ID/WEBHOOK_TOKEN`
- Slack: `slack://TOKEN/CHANNEL`
- Telegram: `tgram://BOT_TOKEN/CHAT_ID`
- Email: `mailto://user:pass@domain`
- Syslog: `syslog://`
- And many more...

See the [Apprise documentation](https://github.com/caronc/apprise#supported-notifications) for a full list of supported services.

## Functions

### Base

Returns the base container with Apprise installed.

```bash
dagger call --mod ./apprise base
```

### Send

Sends a notification via apprise.

Parameters:

- `title` - Notification title
- `body` - Notification body/message
- `service` - Service URL as a secret (use `env:VAR_NAME` to reference environment variables)
