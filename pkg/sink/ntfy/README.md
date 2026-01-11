# Example configuration

```yaml
sinks:
  ntfy:
    type: ntfy
    params:
      url: https://ntfy.sh      # url of the ntfy server
      topic: mytopic            # topic to publish on
      token: tk_******          # your token (optional)
      title: Hello              # title (optional)
      priority: 3               # priority (optional)
```
