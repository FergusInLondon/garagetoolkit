# Garage Toolkit

## Run Locally

```
$ sudo modprobe vcan
$ sudo ip link add dev vcan0 type vcan
$ sudo ip link set up vcan0 mtu 72
```

## CAN Logging

```
$ go build ./cmd/logger
$ ./logger vcan0
```

### Parsing CAN Logs

```
$ go build ./cmd/parser
$ ./parser <path-to-log>
```

### Uploading CAN Logs

**N/A**