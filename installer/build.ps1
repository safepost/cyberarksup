# @echo off

$env:path+=";C:\Program Files (x86)\WiX Toolset v3.11\bin"

candle.exe CASmartSup.wxs
light.exe -ext WixUIExtension .\CASmartSup.wixobj

# https://www.firegiant.com/wix/tutorial/getting-started/
# https://stackoverflow.com/questions/596919/how-to-add-a-ui-to-a-wix-3-installer