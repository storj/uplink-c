module storj.io/uplink-c/testsuite

go 1.13

require (
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/stretchr/testify v1.7.0
	storj.io/common v0.0.0-20210916151047-6aaeb34bb916
	storj.io/drpc v0.0.26
	storj.io/storj v0.12.1-0.20210921100200-32cee1e572f6
	storj.io/uplink v1.6.0
)

replace storj.io/uplink-c => ../
