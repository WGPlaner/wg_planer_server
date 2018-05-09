#!/bin/bash

SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"
PROJECT_DIR=${SCRIPT_DIR}/..

GIT_HASH=$(git --git-dir=".git" show --no-patch --pretty="%h")
echo "GIT_HASH = ${GIT_HASH}"

GIT_DATE=$(git --git-dir=".git" show --no-patch --pretty="%ci")
echo "GIT_DATE = ${GIT_DATE}"

RELEASE_DATE=$(date -u +"%Y-%m-%dT%H:%M:%S%z" --date="${GIT_DATE}")
echo "RELEASE_DATE = ${RELEASE_DATE}"

DATE_HASH=$(date -u +"%Y-%m-%d_%H-%M")
echo "DATE_HASH = ${DATE_HASH}"

VERSION_NAME="${TRAVIS_GO_VERSION}_${DATE_HASH}_git-${GIT_HASH}"
echo "VERSION_NAME = ${VERSION_NAME}"

FILENAME="wg_planer_${VERSION_NAME}.tar.gz"

cd ${PROJECT_DIR}
tar -zcvf wg_planer_server.tar.gz build
cp wg_planer_server.tar.gz ${FILENAME}
cd ${SCRIPT_DIR}

cat > "${SCRIPT_DIR}/bintray.json" <<EOF
{
	"package": {
		"name": "wg_planer_api",
		"repo": "wg_planer_server",
		"subject": "bugwelle",
		"website_url": "https://github.com/wgplaner/wg_planer_server",
		"vcs_url": "https://github.com/wgplaner/wg_planer_server.git",
		"licenses": ["MIT"]
	},
	"version": {
		"name": "${VERSION_NAME}",
		"released": "${RELEASE_DATE}",
		"gpgSign": false
	},
	"files":
	[
		{
			"includePattern": "./${FILENAME}",
			"uploadPattern": "${FILENAME}"
		}
	],
	"publish": true
}
EOF

echo "Config file:"
cat ${SCRIPT_DIR}/bintray.json
