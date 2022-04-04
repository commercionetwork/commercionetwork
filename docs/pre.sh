#!/usr/bin/env bash

mkdir -p docs/modules

for D in ./x/*; do
  if [ -d "${D}" ]; then
    CUR_MOD="docs/modules/$(echo $D | awk -F/ '{print $NF}')"
    N_MOD="$(echo $D | awk -F/ '{print $NF}')"
    rm -rf $CUR_MOD
    if [ -d $D/spec ]; then
      echo "[OK] Coping $N_MOD spec"
      mkdir -p $CUR_MOD && cp -r $D/spec/* "$_"
    fi
  fi
done

cat ./x/README.md | \
  sed 's/\.\/x/\/modules/g' | \
  sed 's/spec\/README.md//g' | \
  sed 's/\.\.\/docs\/building-modules\/README\.md/\/building-modules\/intro\.html/g' > ./docs/modules/README.md