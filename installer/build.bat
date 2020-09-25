@echo off

candle.exe Product.wxs
light.exe -ext WixUIExtension .\Product.wixobj

rem https://www.firegiant.com/wix/tutorial/getting-started/
rem https://stackoverflow.com/questions/596919/how-to-add-a-ui-to-a-wix-3-installer