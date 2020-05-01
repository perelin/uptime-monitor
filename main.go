package main

/*
https://pkg.go.dev/github.com/aws/aws-sdk-go/service/cloudwatch?tab=doc
https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatch/#CloudWatch.PutMetricData

2do later: get instance-id automatically with
wget -q -O - http://instance-data/latest/meta-data/instance-id
https://stackoverflow.com/questions/625644/how-to-get-the-instance-id-from-within-an-ec2-instance

*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func main() {

	var basetime int64

	cliParams := os.Args

	if len(cliParams) == 2 && cliParams[1] == "reset" {
		basetime = setBasetime()
	} else {
		basetime = getBasetime()
	}

	uptime := getUptime(basetime)

	sendUptimeToCloudWatch(float64(uptime))

}

func getBasetime() int64 {
	basetime, err := readBasetime()
	if err != nil {
		fmt.Println(err)
		basetime = setBasetime()
	}
	return basetime
}

func getUptime(basetime int64) int64 {
	uptime := time.Now().Unix() - basetime
	return uptime
}

func setBasetime() int64 {
	nowUnix := time.Now().Unix()
	nowString := strconv.FormatInt(nowUnix, 10)

	basetime := []byte(nowString)
	err := ioutil.WriteFile("/tmp/uptime-monitor-basetime", basetime, 0644)
	if err != nil {
		panic(err)
	}
	return nowUnix
}

func readBasetime() (int64, error) {
	basetimeByte, err := ioutil.ReadFile("/tmp/uptime-monitor-basetime")
	if err != nil {
		return 0, err
	}

	basetimeInt64, err := strconv.ParseInt(string(basetimeByte), 10, 64)
	if err != nil {
		return 0, err
	}

	return basetimeInt64, nil
}

func sendUptimeToCloudWatch(uptime float64) {

	// credentials and rights are managed by an IAM role on the EC2 instance

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	cloudWatchClient := cloudwatch.New(sess)

	var dimensionInstance cloudwatch.Dimension
	dimensionInstance.SetName("InstanceId")
	dimensionInstance.SetValue("i-07570a5f2b15fa22c")

	var data cloudwatch.MetricDatum
	data.SetMetricName("CustomUptimeSeconds")
	data.SetUnit("Seconds")
	data.SetValue(uptime)
	data.SetDimensions([]*cloudwatch.Dimension{&dimensionInstance})

	var input cloudwatch.PutMetricDataInput
	input.SetNamespace("EC2")
	input.SetMetricData([]*cloudwatch.MetricDatum{&data})

	_, err = cloudWatchClient.PutMetricData(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now()

	uptimeDuration, _ := time.ParseDuration(strconv.Itoa(int(uptime)) + "s")
	fmt.Println(
		now.Format("2006-01-02 15:04:05") +
			" - uptime send to cloud watch: " +
			strconv.Itoa(int(uptime)) +
			" -> " +
			uptimeDuration.String())

}

func getSystemUptimeSeconds() (int64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return info.Uptime, nil
}
