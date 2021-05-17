root=$(dirname "$(dirname "$(realpath "$0")")")
filename="food_ordering_backend"

if [[ $(uname) == 'Windows' ]]; then
  filename="food_ordering_backend.exe"
fi

swag init && go build && "$root"/"$filename"
