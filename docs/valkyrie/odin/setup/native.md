# Native Mode

Native mode in Odin uses the host system itself as both the server and worker for code execution. This is the simplest configuration but should be used with caution in production environments.

## Setup

1. Ensure Odin is installed on your system
2. Configure the following environment variables:
   ```bash
   ODIN_MODE=native
   ODIN_WORKER_COUNT=<number_of_workers>  # Optional, defaults to CPU core count
   ODIN_SERVER_PORT=<port_number>         # Optional, defaults to 3000
   ```

## Security Considerations

- Native mode executes code directly on your host system
- Recommended for development and testing environments only
- Consider using ns-jail or container modes for untrusted code execution
- Implement appropriate access controls and rate limiting

## Usage

Start the Odin server:
```bash
odin server start
```

The server will automatically configure workers based on your system's resources.

## Monitoring

- Worker status: `odin worker status`
- Server logs: `odin server logs`
- System resource usage: `odin stats`

## Best Practices

1. Set appropriate resource limits in your configuration
2. Monitor system resource usage regularly
3. Keep Odin updated to the latest version
4. Implement proper error handling in your applications