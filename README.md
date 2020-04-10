# uptime-monitor

Script used to send a relative uptime value of an AWS EC2 instance to AWS cloudwatch. Cloudwatch then can trigger action based on uptime (like stop, reboot, etc).

see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/UsingAlarmActions.html

## create the binary
$ go build .

## install binary
$ mv ./uptime-monitor /usr/local/bin

## usage

### send relative uptime
$ uptime-monitor

### reset uptime
$ uptime-monitor reset

## 2do 
* refactor to parameterize ec2 instance ID 