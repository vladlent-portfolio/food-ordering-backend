cd ../controllers || exit 1
go test ./...
cd ../e2e || exit 1
# disable parallel test execution
go test ./... -p 1
