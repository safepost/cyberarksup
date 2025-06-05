# @echo off

$env:path+=";C:\Program Files (x86)\WiX Toolset v3.14\bin"

candle.exe CyberarkSup.wxs
light.exe -ext WixUIExtension .\CyberarkSup.wixobj

# https://www.firegiant.com/wix/tutorial/getting-started/
# https://stackoverflow.com/questions/596919/how-to-add-a-ui-to-a-wix-3-installer

# Download Wix here
# https://github.com/wixtoolset/wix3/releases