#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/barkisd/${BINARY:-barkisd}
ID=${ID:-0}
LOG=${LOG:-barkisd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'barkisd' E.g.: -e BINARY=barkisd_my_test_version"
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
export BARKISDHOME="/barkisd/node${ID}/barkisd"

if [ -d "`dirname ${BARKISDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$BARKISDHOME" "$@" | tee "${BARKISDHOME}/${LOG}"
else
  "$BINARY" --home "$BARKISDHOME" "$@"
fi

chmod 777 -R /barkisd

