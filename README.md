# Command Monitor

Command Monitor is a web application that allows you to monitor the output of a command and check its response based on various match criteria. The application provides different endpoints to fetch the result in JSON, XML, or plain text formats. Additionally, you can configure the command and match criteria through a web interface.

## Configuration

The configuration is stored in a JSON file named `config.json`. Below is an example of the configuration file:

```json
{
    "command": "echo Hello, World!",
    "match_type": "exact",
    "match_value": "Hello, World!",
    "port": "8080"
}
```
### Configuration Fields
- command: The command to be executed.
- match_type: The type of match to perform on the command output. It can be exact, regex, or integer.
- match_value: The value against which the command output will be matched.
- port: The HTTP port on which the HTTP server will listen.

## Endpoints
### /check/json
- Returns a JSON payload with the result of the command execution.

- Method: GET
- Response:
   ```json
   {
      "success": true,
      "output": "Hello, World!",
      "command": "echo Hello, World!",
      "match_type": "exact",
      "match_value": "Hello, World!"
   }
   ```
### /check/xml
- Returns an XML payload with the result of the command execution.
- Method: GET
- Response:
```xml
<Result>
    <Success>true</Success>
    <Output>Hello, World!</Output>
    <Command>echo Hello, World!</Command>
    <MatchType>exact</MatchType>
    <MatchValue>Hello, World!</MatchValue>
</Result>
```

### /check/status
- Returns a plain text value of success or failed based on the match result.
- Method: GET
- Response:
```
success
```
### /configure
- Provides a web interface to update the configuration values.

### Reload Systemd to Apply the Service
```bash
sudo systemctl daemon-reload
```

### Start the Service:
```bash
sudo systemctl start command-monitor
```

### Enable the Service to Start on Boot:

```bash
sudo systemctl enable command-monitor
```

### Check the Status of the Service:
```bash
sudo systemctl status command-monitor
```