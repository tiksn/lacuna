$ErrorActionPreference = 'Stop';

$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

Get-ChocolateyUnzip -FileFullPath "$toolsDir\winx86.zip" -Destination $toolsDir -FileFullPath64 "$toolsDir\winx64.zip"

Remove-Item ($toolsDir + '\*.' + 'zip')
