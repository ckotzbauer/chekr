module github.com/ckotzbauer/chekr

go 1.16

require (
	github.com/ddelizia/channelify v0.0.1
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go/v12 v12.0.0
)

replace github.com/ddelizia/channelify => github.com/ckotzbauer/channelify v0.0.2-0.20210225173251-a871a15779a0
