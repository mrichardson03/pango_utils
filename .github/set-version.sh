#!/usr/bin/env bash

SCRIPT_BASE="$(cd "$( dirname "$0")" && pwd )"
ROOT=${SCRIPT_BASE}/..

# Exit immediatly if any command exits with a non-zero status
set -e

# Usage
print_usage() {
    echo "Set the app/add-on version"
    echo ""
    echo "Usage:"
    echo "  set-version.sh <new-version>"
    echo ""
}

# if less than one arguments supplied, display usage
if [  $# -lt 1 ]
then
    print_usage
    exit 1
fi

# check whether user had supplied -h or --help . If yes display usage
if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    print_usage
    exit 0
fi

NEW_VERSION=$(echo "$1" | sed -e 's/-beta\./.b/' | sed -e 's/-alpha\./.a/')

echo "### Before relocate..."
ls -R

# Relocate downloaded artifacts.
find build/ -iname '*.zip' -exec mv {} build/ \;

# Rename downloaded artifacts.
for file in $(find build -iname '*.zip');
do
    SUFFIX=$(echo ${file} | cut -d '-' -f2 -f3)
    mv ${file} build/pango_utils-${NEW_VERSION}-${SUFFIX}
done

echo "\n### After relocate..."
ls -R

# Set version in README.md
grep -E '^Version: (.+)$' "$ROOT/README.md" > /dev/null
sed -i.bak -E "s/^Version: (.+)$/Version: $NEW_VERSION/" "$ROOT/README.md" && rm "$ROOT/README.md.bak"
