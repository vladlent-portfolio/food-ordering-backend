root=$(dirname "$(dirname "$(realpath "$0")")")
cd "$root" || exit 1
cd ./controllers || exit 1
go test ./... -count 1
cd ../e2e || exit 1
# disable parallel test execution and caching
go test ./... -p 1 -count=1
