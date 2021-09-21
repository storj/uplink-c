module storj.io/uplink-c/testsuite

go 1.13

require (
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/stretchr/testify v1.7.0
	storj.io/common v0.0.0-20210928143209-230bee624465
	storj.io/drpc v0.0.26
	storj.io/storj v0.12.1-0.20210929104150-f52f5931dafc
	storj.io/uplink v1.6.1-0.20210927115829-4da201e4aebb
)

replace storj.io/uplink-c => ../
