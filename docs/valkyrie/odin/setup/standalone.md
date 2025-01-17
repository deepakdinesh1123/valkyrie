# Standalone Mode

Standalone mode combines the server and worker functionality into a single process, suitable for simpler development and test environments. It also runs an embedded postgres as part of odin itself.

## Features

- Reduced operational complexity
- Suitable for development and testing
- Lower resource overhead
- embedded postgres server

## Usage

Start the standalone server and worker:
```bash
odin standalone
```

## Flags
The `odin standalone` command supports the following flags:

| Flag          | Shorthand | Description                                             | Default Value |
|---------------|-----------|---------------------------------------------------------|---------------|
| `--clean-db`  | `-c`      | Delete existing DB data (deletes both schemas and data) | `false`       |

## Limitations

1. Limited scalability compared to distributed modes
2. Single point of failure
3. Resource constraints of single machine

## Best Practices

1. Use for development or testing
2. Consider other modes for production deployments



