print_opts() {
  echo "options:"
  echo "    darwin:   Build exec file for mac"
  echo "    linux:    Build exec file for linux"
  exit 1
}

if [ $# -le 0 ]; then
  print_opts
fi

if [ $1 == darwin ]; then
    GOOS='darwin' go build -o ./bin/darwin/trs
elif [ $1 == linux ]; then
    GOOS='linux' go build -o ./bin/linux/trs
# 不能编译 net、os 等包到 so 文件中
# elif [ $1 == android ]; then
#     CGO_ENABLED=1 \
#     GOOS=android \
#     GOARCH=amd64 \
#     CC=~/sdk/android/android-ndk-r20b/toolchains/llvm/prebuilt/darwin-x86_64/bin/x86_64-linux-android29-clang++ \
#     go build -buildmode=c-shared -o ./bin/android/libtrs.so
else
    print_opts
fi