# Standalone Mode

Standalone mode combines the server and worker functionality into a single process, suitable for simpler deployments or development environments.

## Configuration

Set the following environment variables:

```bash
ODIN_MODE=standalone
ODIN_PORT=<port_number>
ODIN_MAX_CONCURRENT_EXECUTIONS=<number>
```

## Features

- Single process architecture
- Reduced operational complexity
- Suitable for development and testing
- Lower resource overhead

## Resource Management

Configure execution limits:
```bash
ODIN_EXECUTION_TIMEOUT=<seconds>
ODIN_MEMORY_LIMIT=<bytes>
ODIN_CPU_LIMIT=<percentage>
```

## Usage

Start the standalone server:
```bash
odin standalone start
```

## Monitoring

Monitor the server status:
```bash
odin status
odin metrics
```

## Limitations

1. Limited scalability compared to distributed modes
2. No built-in load balancing
3. Single point of failure
4. Resource constraints of single machine

## Best Practices

1. Set appropriate timeouts and resource limits
2. Monitor system resource usage
3. Implement proper error handling
4. Use for development or low-traffic environments
5. Consider other modes for production deployments

## Troubleshooting

1. Check server status:
   ```bash
   odin status
   ```

2. View logs:
   ```bash
   odin logs
   ```

3. Monitor resource usage:
   ```bash
   odin metrics
   ```