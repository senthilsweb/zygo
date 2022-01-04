#!/bin/bash

OS="$(uname -s)"
GIT_COMMIT=$(git rev-parse --short HEAD)
TAG=$(git describe --exact-match --abbrev=0 --tags ${COMMIT} 2> /dev/null || true)
TAG=${TAG:="dev"}
DATE=$(date +'%Y-%m-%d')
APP_NAME='zygo'
BIN_NAME=$APP_NAME'_'
DIST='./dist/'
DEST='./.temp/'
#INSTALL_BUNDLE=$APP_NAME'_'$TAG'_'$GIT_COMMIT'_'$DATE'.tar.gz'
INSTALL_BUNDLE=$APP_NAME'.tar.gz'
echo $INSTALL_BUNDLE
echo 'Operating System = ['$OS']'

echo "Building binaries"
echo Git commit: $GIT_COMMIT Version: $TAG Build date: $DATE


go generate

# MAC
export GOARCH="amd64"
export GOOS="darwin"
export CGO_ENABLED=1
#go build -ldflags "-X https://github.com/senthilsweb/server.varman.go/cmd.GitCommit=$GIT_COMMIT -X https://github.com/senthilsweb/server.varman.go/cmd.Version=$TAG -X https://github.com/senthilsweb/server.varman.go/cmd.BuildDate=$DATE" -o $DEST$BIN_NAME'mac_amd64' -v .

#LINUX
export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0
#go build -ldflags "-X https://github.com/senthilsweb/server.varman.go/cmd.GitCommit=$GIT_COMMIT -X https://github.com/senthilsweb/server.varman.go/cmd.Version=$TAG -X https://github.com/senthilsweb/server.varman.go/cmd.BuildDate=$DATE" -o $DEST$BIN_NAME'linux_amd64' -v


export GOARCH="arm"
export GOOS="linux"
export CGO_ENABLED=0   
export GOARM=5 
go build -ldflags "-X https://github.com/senthilsweb/server.varman.go/cmd.GitCommit=$GIT_COMMIT -X https://github.com/senthilsweb/server.varman.go/cmd.Version=$TAG -X https://github.com/senthilsweb/server.varman.go/cmd.BuildDate=$DATE" -o $DEST$BIN_NAME'linux_arm' -v

#WINDOWS
export GOARCH="386"
export GOOS="windows"
export CGO_ENABLED=0
#go build -ldflags "-X https://github.com/senthilsweb/server.varman.go/cmd.GitCommit=$GIT_COMMIT -X https://github.com/senthilsweb/server.varman.go/cmd.Version=$TAG -X https://github.com/senthilsweb/server.varman.go/cmd.BuildDate=$DATE" -o $DEST$BIN_NAME'windows_386.exe' -v

export GOARCH="amd64"
export GOOS="windows"
export CGO_ENABLED=0
#go build -ldflags "-X https://github.com/senthilsweb/server.varman.go/cmd.GitCommit=$GIT_COMMIT -X https://github.com/senthilsweb/server.varman.go/cmd.Version=$TAG -X https://github.com/senthilsweb/server.varman.go/cmd.BuildDate=$DATE" -o $DEST$BIN_NAME'windows_amd64.exe' -v

#cp -rf ./configs $DEST
cp app.service.txt $DEST$APP_NAME'.service.txt'
#cp config.yml $DEST'config.yml'
cp README.md $DEST'README.md'
cp app.install.md $DEST$APP_NAME'.install.md'
cp install.sh $DEST'install.sh'

if [ "$(uname)" == "Darwin" ]; then
    #For Mac
    sed -i "" 's/{{app}}/'$APP_NAME'/g' $DEST$APP_NAME'.service.txt'
    #sed -i "" 's/{{app}}/'$APP_NAME'/g' $DEST$APP_NAME'.env'
    sed -i "" 's/{{app}}/'$APP_NAME'/g' $DEST'install.sh'
    #sed -i "" 's/{{app}}/'$APP_NAME'/g' $DEST'config.json'
    #sed -i "" 's/{{app}}/'$APP_NAME'/g' $DEST'README.md'
    sed -i "" 's/{{app}}/'$APP_NAME'/g' $DEST$APP_NAME'.install.md'
    sed -i "" 's/{{tar_file}}/'$INSTALL_BUNDLE'/g' $DEST$APP_NAME'.install.md'
else
    #For Linux
    sed -i 's/{{app}}/'$APP_NAME'/g' $DEST$APP_NAME'.service.txt'
    #sed -i 's/{{app}}/'$APP_NAME'/g' $DEST$APP_NAME'.env'
    sed -i 's/{{app}}/'$APP_NAME'/g' $DEST'install.sh'
    #sed -i 's/{{app}}/'$APP_NAME'/g' $DEST'config.json'
    #sed -i 's/{{app}}/'$APP_NAME'/g' $DEST'README.md'
    sed -i 's/{{app}}/'$APP_NAME'/g' $DEST$APP_NAME'.install.md'
    sed -i 's/{{tar_file}}/'$INSTALL_BUNDLE'/g' $DEST$APP_NAME'.install.md'
fi
tar -czvf $DIST$INSTALL_BUNDLE -C $DEST .
#rm -rf $DEST
echo "Build complete"