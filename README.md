# uptime-monitor

Script used to send a relative uptime value to AWS cloudwatch. Cloudwatch then can trigger action based on uptime (like shutdown).

## create the binary
$ go build .

## install binary
$ mv ./uptime-monitor /usr/local/bin

## usage

### send relative uptime
$ uptime-monitor

### reset uptime
$ uptime-monitor reset