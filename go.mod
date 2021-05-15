module github.com/ckotzbauer/chekr

go 1.16

require (
	github.com/ddelizia/channelify v0.0.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.24.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/text v0.3.5 // indirect
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
)

replace github.com/ddelizia/channelify => github.com/ckotzbauer/channelify v0.0.2-0.20210225173251-a871a15779a0
