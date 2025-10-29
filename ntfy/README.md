# ntfy

A Dagger module for sending notifications via [ntfy](https://ntfy.sh/).

## Usage

### Basic notification

```bash
dagger call --mod ./ntfy send --topic=phil_alerts --message="Hello from Dagger"
```

### With all optional parameters

```bash
dagger call --mod ./ntfy send \
  --topic=phil_alerts \
  --message="Remote access detected" \
  --title="Unauthorized access detected" \
  --priority=urgent \
  --tags="warning,skull"
```

### Using a custom ntfy server

```bash
dagger call --mod ./ntfy send \
  --topic=my_topic \
  --message="Hello" \
  --server="https://ntfy.example.com"
```

### With Markdown formatting

```bash
dagger call --mod ./ntfy send \
  --topic=my_topic \
  --message="Look ma, **bold text**, *italics*, [links](https://example.com)" \
  --markdown=true
```

## Functions

### Send

Sends a notification to an ntfy topic.

**Parameters:**

- `topic` (required) - The ntfy topic to send to
- `message` (required) - The message to send
- `title` (optional) - The notification title
- `server` (optional) - The ntfy server URL (default: `https://ntfy.sh`)
- `priority` (optional) - Priority of the notification (e.g., `urgent`, `high`, `default`, `low`, `min`)
- `tags` (optional) - Comma-separated list of tags (e.g., `warning,skull`)
- `markdown` (optional) - Enable Markdown formatting for bold, italics, links, images, code blocks, etc.
