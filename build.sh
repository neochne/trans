print_opts() {
  echo "options:"
  echo "    darwin: Build for mac"
  echo "    linux:  Build for linux"
  exit 1
}

if [ $# -le 0 ]; then
  print_opts
fi

if [ $1 == darwin ]; then
    go build -o ./bin/darwin/trans
elif [ $1 == linux ]; then
    GOOS='linux' go build -o ./bin/linux/trans
else
    print_opts
fi