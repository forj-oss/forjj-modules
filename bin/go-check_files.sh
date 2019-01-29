#!/bin/bash
#

if [[ "$BE_PROJECT" = "" ]]
then
    echo "Missing project environment. Load it with build-env"
    exit 1
fi
FILES="$(find . -wholename ./vendor -prune -o -wholename ./.glide -prune -o -name "*.go" -exec grep \"$BE_PROJECT/ {} \; -print)"
if [[ "$(echo $FILES | wc -l)" -ne 0 ]]
then
    echo "A GO module requires self reference to be exported. (relative path is not accepted)
    
List of files with wrong import:
$FILES"
    exit 1
fi