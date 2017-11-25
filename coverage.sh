#!/usr/bin/env bash

PACKAGES=$(go list ./... | grep -v -e vendor -e restapi)

go get -u github.com/wadey/gocovmerge

# Unit tests
for PKG in $PACKAGES; do
	go test -cover -coverprofile $GOPATH/src/$PKG/coverage.out $PKG
done;

# Integration tests
go test -c github.com/wgplaner/wg_planer_server/integrations -coverpkg $(echo $PACKAGES | tr ' ' ',') -o integrations.cover.test
./integrations.cover.test -test.coverprofile=integration.coverage.out

echo "mode: set" > coverage.all
for PKG in $PACKAGES; do
	egrep "$PKG[^/]*\.go" integration.coverage.out
	egrep "$PKG/[^/]*\.go" integration.coverage.out > int.coverage.out;

	# TODO
	#gocovmerge $GOPATH/src/$PKG/coverage.out int.coverage.out > pkg.coverage.out;
	cat $GOPATH/src/$PKG/coverage.out int.coverage.out > pkg.coverage.out

	grep -h -v "^mode:" pkg.coverage.out >> coverage.all;
	mv pkg.coverage.out $GOPATH/src/$PKG/coverage.out;
	rm int.coverage.out;
done;
