Properties {
    $version = "0.0.1"
}

Task PublishChocolateyPackage -Depends PackChocolateyPackage {
    Exec { choco push $script:chocoNupkg }
}

Task PublishLinuxRpmPackages -Depends PackLinuxPackage {
    #
}

Task PackLinuxPackage -Depends BuildLinux64, BuildRhel64 {
}

Task PackChocolateyPackage -Depends ZipBuildArtifacts {
    Copy-Item -Path .\Chocolatey\tools\chocolateyInstall.ps1 -Destination $script:chocoTools
    Copy-Item -Path .\Chocolatey\tools\LICENSE.txt -Destination $script:chocoLegal
    Copy-Item -Path .\Chocolatey\tools\VERIFICATION.txt -Destination $script:chocoLegal

    $verificationFilePath = Join-Path -Path $script:chocoLegal -ChildPath VERIFICATION.txt

    $zipX64Hash = Get-FileHash -Path $script:artifactsZipX64
    $zipX86Hash = Get-FileHash -Path $script:artifactsZipX86

    function AddAllHashes ($Path, $Depth = 0) {
        $items = Get-ChildItem -Path $Path -File

        foreach($item in $items){
            $fileName = $item.Name
            $hash = Get-FileHash -Path $item
            $algorithm = $hash.Algorithm
            $hash = $hash.Hash
            Add-Content -Path $verificationFilePath -Value "Archived file $fileName has Hash: $algorithm $hash"
        }
    }
    
    Add-Content -Path $verificationFilePath -Value ("File Hash: winx64.zip - " + $zipX64Hash.Algorithm + " - " + $zipX64Hash.Hash)
    AddAllHashes($script:publishWinx64Folder)
    Add-Content -Path $verificationFilePath -Value ""
    Add-Content -Path $verificationFilePath -Value ("File Hash: winx86.zip - " + $zipX86Hash.Algorithm + " - " + $zipX86Hash.Hash)
    AddAllHashes($script:publishWinx86Folder)
    Add-Content -Path $verificationFilePath -Value ""
    Add-Content -Path $verificationFilePath -Value ""

    $repoStatus = Get-RepositoryStatus
    $commitHash = $repoStatus.CurrentCommit
    Add-Content -Path $verificationFilePath -Value "Git commit hash: $commitHash"

    $chocoNuspec = Join-Path -Path $script:chocolateyPublishFolderFolder -ChildPath lacuna.nuspec
    Copy-Item -Path ".\Chocolatey\lacuna.nuspec" -Destination $chocoNuspec
    Exec { choco pack $chocoNuspec --version $version --outputdirectory $script:trashFolder version=$version }
    $script:chocoNupkg = Join-Path -Path $script:trashFolder -ChildPath "lacuna.$version.nupkg"
}

Task ZipBuildArtifacts -Depends BuildWinx64, BuildWinx86 {
    $script:artifactsZipX64 = Join-Path -Path $script:chocoTools -ChildPath "winx64.zip"
    $script:artifactsZipX86 = Join-Path -Path $script:chocoTools -ChildPath "winx86.zip"

    Compress-Archive -Path "$script:publishWinx64Folder\*" -CompressionLevel Optimal -DestinationPath $script:artifactsZipX64
    Compress-Archive -Path "$script:publishWinx86Folder\*" -CompressionLevel Optimal -DestinationPath $script:artifactsZipX86
}

Task Build -Depends BuildWinx64, BuildWinx86, BuildLinux64

Task BuildWinx64 -Depends PreBuild {
    $script:publishWinx64Folder = Join-Path -Path $script:publishFolder -ChildPath "winx64"
    $outputFile = Join-Path -Path $script:publishWinx64Folder -ChildPath "lacuna.exe"

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile $script:srcFolder }
}

Task BuildWinx86 -Depends PreBuild {
    $script:publishWinx86Folder = Join-Path -Path $script:publishFolder -ChildPath "winx86"
    $outputFile = Join-Path -Path $script:publishWinx86Folder -ChildPath "lacuna.exe"
    
    $env:GOOS = "windows"
    $env:GOARCH = "386"
    Exec { go build -o $outputFile $script:srcFolder }
}

Task BuildLinux64 -Depends PreBuild {
    $script:publishLinux64Folder = Join-Path -Path $script:publishFolder -ChildPath "linux64"
    $outputFile = Join-Path -Path $script:publishLinux64Folder -ChildPath "lacuna"

    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile $script:srcFolder }
}

Task PreBuild -Depends Init, Clean, Format, InstallPackages {
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
    Exec { go get "github.com/docker/distribution/reference" }
}

Task Format -Depends Clean {
    Exec { go fmt $script:srcFolder }
}

Task Clean -Depends Init {
    Exec { go clean $script:srcFolder }
}

Task Init {
    $date = Get-Date
    $ticks = $date.Ticks
    $trashFolder = Join-Path -Path . -ChildPath ".trash"
    $script:trashFolder = Join-Path -Path $trashFolder -ChildPath $ticks.ToString("D19")
    New-Item -Path $script:trashFolder -ItemType Directory | Out-Null
    $script:trashFolder = Resolve-Path -Path $script:trashFolder
    $script:srcFolder = Resolve-Path -Path ".\src\" -Relative
}
 