# NS Jail Mode

NS Jail mode provides enhanced security for code execution by utilizing Linux namespaces and resource isolation through nsjail.

## Prerequisites

- nsjail installed on your system
- Linux kernel with namespace support
- Proper system permissions

## Configuration

Set the following environment variables:

```bash
ODIN_MODE=nsjail
ODIN_NSJAIL_CONFIG=/path/to/nsjail.cfg
ODIN_NSJAIL_TIME_LIMIT=<seconds>
ODIN_NSJAIL_MEMORY_LIMIT=<bytes>
```

## Security Features

1. Process Isolation
   - Separate PID namespace
   - Custom mount namespace
   - Network namespace isolation

2. Resource Limits
   - CPU time restrictions
   - Memory usage caps
   - File size quotas
   - Maximum process count

3. Filesystem Isolation
   - Read-only root
   - Restricted /proc mount
   - Temporary scratch space

## Example Configuration

Create a nsjail configuration file:

```bash
name: "odin-nsjail"
mode: ONCE
rlimit_as: 2048
rlimit_cpu: 1000
rlimit_fsize: 1024
mount {
    src: "/bin"
    dst: "/bin"
    is_bind: true
    rw: false
}
mount {
    src: "/lib"
    dst: "/lib"
    is_bind: true
    rw: false
}
```

## Monitoring and Logging

- System logs: Check syslog for nsjail-related entries
- Resource usage: Monitor through cgroup statistics
- Execution logs: Available in specified log directory

## Troubleshooting

1. Check nsjail logs:
   ```bash
   journalctl | grep nsjail
   ```

2. Verify namespace isolation:
   ```bash
   ls -l /proc/self/ns/
   ```

3. Test configuration:
   ```bash
   nsjail --config /path/to/nsjail.cfg --test
   ```
