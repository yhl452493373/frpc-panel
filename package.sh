#!/bin/bash
version=$(cat cmd/frpc-panel/cmd.go | grep 'const version' | egrep -o '[0-9.]+')
cd ./release || exit
rm -rf *.zip
list=$(ls frpc-panel-*)
echo "$list"
for binFile in $list
  do
    tmpFile=frpc-panel
    newBinFile=$binFile
    if echo "$binFile" | grep  -q -E "\.exe";then
      tmpFile=frpc-panel.exe
      newBinFile=${newBinFile%%.exe*}
    fi
    cp "$binFile" "$tmpFile"
    zip -r "$newBinFile-$version".zip "$tmpFile" frpc-panel.toml assets -x "*.git*" "*.idea*" "*.DS_Store" "*.contentFlavour"
    rm -rf "$binFile" "$tmpFile"
  done
  rm -rf frpc-panel.toml assets
