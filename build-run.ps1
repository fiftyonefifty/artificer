Param
(
$version
)
Write-Host $version


Write-Host "Script:" $PSCommandPath
Write-Host "Path:" $PSScriptRoot

Set-Location -Path  $PSScriptRoot

function Run-Script
{
	param([string]$script)
	$ScriptPath = "$PSScriptRoot\$script.ps1"
	& $ScriptPath
}

$Time = [System.Diagnostics.Stopwatch]::StartNew()

function PrintElapsedTime {
    Log $([string]::Format("Elapsed time: {0}.{1}", $Time.Elapsed.Seconds, $Time.Elapsed.Milliseconds))
}

function Log {
    Param ([string] $s)
    Write-Output "###### $s"
}

function Check {
    Param ([string] $s)
    if ($LASTEXITCODE -ne 0) { 
        Log "Failed: $s"
        throw "Error case -- see failed step"
    }
}

$DockerOS = docker version -f "{{ .Server.Os }}"
write-host DockerOs:$DockerOS 
$BaseBuildImageName = "artificer"
$Dockerfile = "build/package/Dockerfile"

PrintElapsedTime

Log "Build application image"
docker build --no-cache --pull -t $BaseBuildImageName -f $PSScriptRoot/$Dockerfile  .
PrintElapsedTime
Check "docker build (application)"

