#!/bin/bash
#

if [[ "$BE_PROJECT" = "" ]]
then
    echo "Missing project environment. Load it with build-env"
    exit 1
fi

if [[ "$1" = "--fix-it" ]] && [[ "$GO_MODULE_REF" != "" ]]
then
    find . -wholename ./vendor -prune -o -wholename ./.glide -prune -o -name "*.go" -exec grep \"$BE_PROJECT/ {} \; -print | while read LINE
    do
        if [[ "$LINE" =~ ^\./ ]]
        then
            echo "$LINE"
            sed -i 's|"'"$BE_PROJECT"'/|"'"$GO_MODULE_REF"'/|g' $LINE
        fi
    done
    exit
fi

FILES="$(find . -wholename ./vendor -prune -o -wholename ./.glide -prune -o -name "*.go" -exec grep \"$BE_PROJECT/ {} \; -print)"
if [[ "$FILES" != "" ]]
then
    echo "A GO module requires self reference to be exported. (relative path is not accepted)
    
List of files with wrong import:
$FILES"
    exit 1
fi
echo "No issue with "