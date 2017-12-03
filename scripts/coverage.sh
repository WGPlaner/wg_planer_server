#!/usr/bin/env bash

# Stupid check if we're in the correct directory
if [ ! -f swagger.yml ]
then
    echo [Error] Execute this script from the wg_planer_server directory!
fi

# Packages to test
PACKAGES=$(go list ./... | grep -v -e vendor -e restapi -e integrations)

# Unit tests
echo "mode: count" > unit.coverage.out
for PKG in ${PACKAGES}; do
	echo "" > package.coverage.out
	go test -covermode count -coverprofile package.coverage.out ${PKG}
	grep -h -v -e "^mode:" -e '^$' package.coverage.out >> unit.coverage.out
done

# Integration tests
if [[ -z "${WGPLANER_ROOT}" ]]; then
  echo "Environment variable WGPLANER_ROOT must be set for integration tests!"
  exit 1
fi

go test -c github.com/wgplaner/wg_planer_server/integrations -covermode count -coverpkg $(echo ${PACKAGES} | tr ' ' ',') -o integrations.cover.test
./integrations.cover.test -test.coverprofile=integration.coverage.out
