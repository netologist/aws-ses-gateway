package main

import (
	"github.com/netologist/aws-ses-gateway/internal"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("AWS SES Gateway")
	internal.StartServer()
}
