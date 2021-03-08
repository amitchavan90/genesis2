#!/usr/bin/env bash

# get minor version number from passed argument if provided OTHERWISE get commit count
c=$1
if [[ -z $c ]]
then
    c=$(git rev-list --count HEAD)
fi

# set version number
v=0.1.$c
echo "building apk version" $v "(" $c ")"

# build
flutter clean
flutter build apk --build-number=$c --build-name=$v

# rename
mv build/app/outputs/apk/release/app-release.apk build/app/outputs/apk/release/genesis-$v.apk

# copy to dist downloads folder
mkdir -p ./build/
cp build/app/outputs/apk/release/genesis-$v.apk ./build/genesis-$v.apk

# update fieldapp version in backend server
l='const fieldappVersion ='
sed -i "/$l/s/.*/$l \"$v\"/" ../server/cmd/platform/main.go
