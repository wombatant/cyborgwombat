echo "package main; const cyborgbear_version = \"$2\";" > version.go
go fmt

liccor

# rebuild the example to get rid of the license headers added by liccor
make -C example

git commit -m "Updated version flag to ${2}." version.go

git tag ${1}-${2}
