#!/usr/bin/env bash

build_dir='build'

package_name='dev-router'

platforms=("windows/amd64" "windows/386" "darwin/amd64" "darwin/386" "linux/amd64" "linux/386")

set -e

echo "cleaning build dir ..."
if [ -d "$build_dir" ]; then
    rm -rf ./$build_dir/*
else
    mkdir -p ./$build_dir
fi

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    zipped_name=$package_name'.'$GOOS'_'$GOARCH'.zip'
    output_name=$package_name
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "building $package_name for $platform ..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ./$build_dir/$output_name
    (cd ./$build_dir && zip $zipped_name $output_name)
done

echo "done."
