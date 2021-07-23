root=$(dirname "$(dirname "$(realpath "$0")")")
swagname="swag"
filename="food_ordering_backend"

# Check that user is using bash for Windows.
if [[ "$(expr substr "$(uname -s)" 1 10)" == "MINGW64_NT" ]]; then
  swagname="swag.exe"
  filename="food_ordering_backend.exe"
fi

cd "$root" || exit 1

if [ -x "$(command -v $swagname)" ]; then
  $swagname init
else
  echo "Couldn't find $swagname in PATH."
  echo "Swagger schema won't be updated."
fi

go build && "$root"/"$filename"
