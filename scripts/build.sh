root=$(dirname "$(dirname "$(realpath "$0")")")
os="linux"
arch="amd64"
outputname="food_ordering_api"

swagname="swag"
rundir="$(pwd)"

cd "$root" || exit

# Check that user is using bash for Windows.
if [[ "$(expr substr "$(uname -s)" 1 10)" == "MINGW64_NT" ]]; then
  swagname="swag.exe"
fi

if [ -x "$(command -v $swagname)" ]; then
  $swagname init
else
  echo "Couldn't find $swagname in PATH."
  echo "Swagger schema won't be updated."
fi

GOOS=$os GOARCH=$arch go build -o "$rundir/$outputname"
