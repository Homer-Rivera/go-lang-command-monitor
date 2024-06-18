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