#!/bin/sh

usage() {
	echo "Usage: $0 <schema.json> [additional check-jsonschema args...]"
	exit 2
}

# Ensure at least a schema file is provided
if [ $# -lt 1 ]; then
	usage
fi

SCHEMAFILE="$1"
shift
PARENTDIR="$(cd "$(dirname "$PWD/$SCHEMAFILE")" && pwd)" # `check-jsonschema` attempts to fetch from the `$id` parameter without this; the resolution should be relative to the local file
#echo $PARENTDIR

check-jsonschema --base-uri=file://$PARENTDIR/ --schemafile "$SCHEMAFILE" "$@"
