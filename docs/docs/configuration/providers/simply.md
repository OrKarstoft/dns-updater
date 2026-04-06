---
sidebar_position: 4
---

# Simply.com

The Simply provider uses the [Simply.com API](https://www.simply.com/en/docs/api/) to update your DNS records.

## Configuration

To use the Simply provider, you need an API key and your account name. You can find these in your Simply.com control panel under the API section.

### Example `config.yaml`

```yaml
provider:
  name: simply
  accountName: "YOUR_ACCOUNT_NAME"
  apiKey: "YOUR_API_KEY"
```
