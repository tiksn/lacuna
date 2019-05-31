Properties {
    $version="0.0.1"
}

Task PublishChocolateyPackage -Depends PackChocolateyPackage {
    Exec { choco push $script:chocoNupkg }
}

Task PublishLinuxRpmPackages -Depends PackLinuxPackage {
    #
}

Task PackLinuxPackage -Depends BuildLinux64,BuildRhel64 {
}

Task PackChocolateyPackage -Depends ZipBuildArtifacts {
    Copy-Item -Path .\Chocolatey\tools\chocolateyInstall.ps1 -Destination $script:chocoTools
    Copy-Item -Path .\Chocolatey\tools\LICENSE.txt -Destination $script:chocoLegal
    Copy-Item -Path .\Chocolatey\tools\VERIFICATION.txt -Destination $script:chocoLegal

    $verificationFilePath = Join-Path -Path $script:chocoLegal -ChildPath VERIFICATION.txt

    $zipX64Hash = Get-FileHash -Path $script:artifactsZipX64
    $zipX86Hash = Get-FileHash -Path $script:artifactsZipX86

    Add-Content -Path $verificationFilePath -Value ("File Hash: winx64.zip - " + $zipX64Hash.Algorithm + " - " + $zipX64Hash.Hash)
    Add-Content -Path $verificationFilePath -Value ("File Hash: winx86.zip - " + $zipX86Hash.Algorithm + " - " + $zipX86Hash.Hash)
    
    $chocoNuspec = Join-Path -Path $script:chocolateyPublishFolderFolder -ChildPath lacuna.nuspec
    Copy-Item -Path ".\Chocolatey\lacuna.nuspec" -Destination $chocoNuspec
    Exec { choco pack $chocoNuspec --version $version --outputdirectory $script:trashFolder version=$version }
    $script:chocoNupkg = Join-Path -Path $script:trashFolder -ChildPath "lacuna.$version.nupkg"
}

Task ZipBuildArtifacts -Depends BuildWinx64,BuildWinx86 {
    $script:artifactsZipX64 = Join-Path -Path $script:chocoTools -ChildPath "winx64.zip"
    $script:artifactsZipX86 = Join-Path -Path $script:chocoTools -ChildPath "winx86.zip"

    Compress-Archive -Path "$script:publishWinx64Folder\*" -CompressionLevel Optimal -DestinationPath $script:artifactsZipX64
    Compress-Archive -Path "$script:publishWinx86Folder\*" -CompressionLevel Optimal -DestinationPath $script:artifactsZipX86
}

Task Build -Depends BuildWinx64,BuildWinx86,BuildLinux64

Task BuildWinx64 -Depends PreBuild {
   $script:publishWinx64Folder = Join-Path -Path $script:publishFolder -ChildPath "winx64"
   $outputFile = Join-Path -Path $script:publishWinx64Folder -ChildPath "lacuna.exe"

   Exec { go build -o $outputFile -tags GOOS=windows -tags GOARCH=amd64 ".\src\" }
}

Task BuildWinx86 -Depends PreBuild {
    $script:publishWinx86Folder = Join-Path -Path $script:publishFolder -ChildPath "winx86"
    $outputFile = Join-Path -Path $script:publishWinx86Folder -ChildPath "lacuna.exe"

    Exec { go build -o $outputFile -tags GOOS=windows -tags GOARCH=386 ".\src\" }
}

Task BuildLinux64 -Depends PreBuild {
    $script:publishLinux64Folder = Join-Path -Path $script:publishFolder -ChildPath "linux64"
    $outputFile = Join-Path -Path $script:publishLinux64Folder -ChildPath "lacuna"

    Exec { go build -o $outputFile -tags GOOS=linux -tags GOARCH=amd64 ".\src\" }
}

Task PreBuild -Depends Init,Clean,Format,InstallPackages {
    $script:publishFolder = Join-Path -Path $script:trashFolder -ChildPath "bin"
    $script:chocolateyPublishFolderFolder = Join-Path -Path $script:trashFolder -ChildPath "choco"
    $script:chocoLegal = Join-Path -Path $script:chocolateyPublishFolderFolder -ChildPath "legal"
    $script:chocoTools = Join-Path -Path $script:chocolateyPublishFolderFolder -ChildPath "tools"
    
    New-Item -Path $script:chocoLegal -ItemType Directory | Out-Null
    New-Item -Path $script:chocoTools -ItemType Directory | Out-Null
    New-Item -Path $script:publishFolder -ItemType Directory | Out-Null
}

Task InstallPackages {
    Exec { go get "github.com/urfave/cli" }
    Exec { go get "github.com/moby/buildkit/frontend/dockerfile/parser" }
}

Task Format -Depends Clean {
    Exec { go fmt ".\src\" }
}
Task Clean -Depends Init {
    Exec { go clean ".\src\" }
}

Task Init {
   $date = Get-Date
   $ticks = $date.Ticks
   $trashFolder = Join-Path -Path . -ChildPath ".trash"
   $script:trashFolder = Join-Path -Path $trashFolder -ChildPath $ticks.ToString("D19")
   New-Item -Path $script:trashFolder -ItemType Directory | Out-Null
   $script:trashFolder = Resolve-Path -Path $script:trashFolder
}
 