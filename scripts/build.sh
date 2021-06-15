root=$(dirname "$(dirname "$(realpath "$0")")")
cd "$root" || exit
swag init
GOOS=linux GOARCH=amd64 go build -o f_o_api
