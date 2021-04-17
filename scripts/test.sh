root=$(dirname "$(dirname "$(realpath "$0")")")
cd "$root" || exit 1
cd ./controllers || exit 1
go test ./...
cd ../e2e || exit 1
# disable parallel test execution
go test ./... -p 1
