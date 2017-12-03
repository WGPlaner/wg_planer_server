#!/usr/bin/env bash

# Stupid check if we're in the correct directory
if [ ! -f swagger.yml ]
then
    echo [Error] Execute this script from the wg_planer_server directory!
fi

if [ ! -f swagger-codegen-cli.jar ]
then
	wget -O swagger-codegen-cli.jar https://oss.sonatype.org/content/repositories/releases/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar
fi

rm -rf ./build/android_client
grep -v "readOnly: true$" swagger.yml  > swagger_android.yml
java -jar swagger-codegen-cli.jar generate -i swagger_android.yml -l java --library=okhttp-gson -o build/android_client

# Build the android library
cd ./build/android_client
mvn clean package
