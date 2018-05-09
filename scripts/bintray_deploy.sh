#!/bin/bash

SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"
PROJECT_DIR=${SCRIPT_DIR}/..

GIT_DATE=$(git --git-dir=".git" show --no-patch --pretty="%ci")
echo "GIT_DATE = ${GIT_DATE}"

RELEASE_DATE=$(date -u +"%Y-%m-%dT%H:%M:%S%z" --date="${GIT_DATE}")
echo "RELEASE_DATE = ${RELEASE_DATE}"

DATE_HASH=$(date -u +"%Y-%m-%d_%H-%M")
echo "DATE_HASH = ${DATE_HASH}"

cp wg_planer_server.tar.gz wg_planer_server_${DATE_HASH}.tar.gz

cat > "${SCRIPT_DIR}/bintray.json" <<EOF
{
	"package": {
		"name": "wg_planer_server",
		"repo": "wg_planer_api",
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
			"includePattern": "${PROJECT_DIR}/wg_planer_server_${DATE_HASH}.tar.gz",
			"uploadPattern": "wg_planer_server_${DATE_HASH}.tar.gz"
		}
	],
	"publish": true
}
EOF
