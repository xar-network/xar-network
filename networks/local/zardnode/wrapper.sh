#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/xard/${BINARY:-xard}
ID=${ID:-0}
LOG=${LOG:-xard.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'xard' E.g.: -e BINARY=xard_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export XARDHOME="/xard/node${ID}/xard"

if [ -d "$(dirname "${XARDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${XARDHOME}" "$@" | tee "${XARDHOME}/${LOG}"
else
  "${BINARY}" --home "${XARDHOME}" "$@"
fi
