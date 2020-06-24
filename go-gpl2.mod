module storj.io/uplink-c

go 1.13

require (
	github.com/calebcase/tmpfile v1.0.2-0.20200602150926-3af473ef8439 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/zeebo/errs v1.2.2
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
	storj.io/common v0.0.0-20200611114417-9a3d012fdb62
	storj.io/drpc v0.0.13 // indirect
	storj.io/uplink v1.1.2
)

replace github.com/spacemonkeygo/monkit/v3 => ./internal/replacements/monkit
