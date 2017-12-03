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

rm -rf ./build/api_doc_html
java -jar swagger-codegen-cli.jar generate -i swagger.yml -l html2 -o build/api_doc_html
