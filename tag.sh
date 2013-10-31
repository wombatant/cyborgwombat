echo "package main; const cyborgbear_version = \"$1\";" > version.go
liccor
git tag ${1}-${2}

# rebuild the example to get rid of the license headers added by liccor
make -C example
