#!/bin/bash

VERSION=$1
DESCRIPTION=$2
TOKEN=$3

if [[ ! $VERSION || ! $DESCRIPTION || ! $TOKEN ]] ; then
    echo "$0 VERSION DESCRIPTION GITHUB_TOKEN"
    exit 1
fi

sed -i -e "s/^const version string = \".*\"/const version string = \"${VERSION}\"/" version.go
git add version.go
git commit -m "Bump version number to ${VERSION}"

git push origin master

echo Creating draft version ${VERSION}...
gothub release -u justsocialapps -r holmes -t v${VERSION} -n ${VERSION} -d "${DESCRIPTION}" --draft -s ${TOKEN}

echo Building release files...
GOOS=windows GOARCH=amd64 go build -o holmes_${VERSION}_windows_amd64.exe
GOOS=windows GOARCH=386 go build -o holmes_${VERSION}_windows_386.exe
GOOS=linux GOARCH=386 go build -o holmes_${VERSION}_linux_386
GOOS=linux GOARCH=amd64 go build -o holmes_${VERSION}_linux_amd64

echo Uploading release files...
for i in holmes_${VERSION}_* ; do
    gothub upload -u justsocialapps -r holmes -s ${TOKEN} -t v${VERSION} -f $i -n $i
done

echo
echo "Please review your release at https://github.com/justsocialapps/holmes/releases"
echo -n "Publish release ${VERSION}? (Y/n) "
read publish
echo

if [[ ${publish} == "" || ${publish} == "y" || ${publish} == "Y" ]] ; then
    echo Publishing release ${VERSION}...
    gothub edit -u justsocialapps -r holmes -t v${VERSION} -n ${VERSION} -d "${DESCRIPTION}" -s ${TOKEN}
    published=1
else
    echo "Discarding version bump..."
    git reset HEAD~1
    git checkout version.go
    echo "Deleting draft release..."
    gothub delete -u justsocialapps -r holmes -t v${VERSION} -s ${TOKEN}
fi

rm holmes_${VERSION}_*

if [[ $published ]] ; then
    echo "Done. You'll find your release here: https://github.com/justsocialapps/holmes/releases/tag/v${VERSION}"
else
    echo "Done."
fi
