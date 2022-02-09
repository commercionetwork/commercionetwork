#!/usr/bin/env bash

mkdir -p docs/modules

for D in ./x/*; do
  if [ -d "${D}" ]; then
    rm -rf "docs/modules/$(echo $D | awk -F/ '{print $NF}')"
    mkdir -p "docs/modules/$(echo $D | awk -F/ '{print $NF}')" && cp -r $D/spec/* "$_"
  fi
done

cat ./x/README.md | \
  sed 's/\.\/x/\/modules/g' | \
  sed 's/spec\/README.md//g' | \
  sed 's/\.\.\/docs\/building-modules\/README\.md/\/building-modules\/intro\.html/g' > ./docs/modules/README.md