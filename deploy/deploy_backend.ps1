param(
    [switch]$SkipBuild
)

$ErrorActionPreference = 'Stop'
$deployDir = $PSScriptRoot
$root = Split-Path -Parent $deployDir
$targetFile = Join-Path $deployDir 'target.local.json'
if (-not (Test-Path $targetFile)) {
    throw "Copy deploy/target.example.json to deploy/target.local.json and fill in the server details."
}
$target = Get-Content $targetFile -Raw | ConvertFrom-Json
if ($target.host -notmatch '^[A-Za-z0-9.-]+$' -or $target.host -eq 'SERVER_IP' -or
    $target.user -notmatch '^[A-Za-z_][A-Za-z0-9_-]*$' -or
    $target.remoteRoot -notmatch '^/[A-Za-z0-9._/-]+$' -or $target.remoteRoot -eq '/' -or
    $target.remoteRoot -match '(^|/)\.\.?(/|$)' -or
    $target.serviceUser -notmatch '^[A-Za-z_][A-Za-z0-9_-]*$' -or
    [int]$target.sshPort -lt 1 -or [int]$target.sshPort -gt 65535) {
    throw 'Invalid deploy/target.local.json.'
}

function Invoke-Checked([string]$command, [string[]]$arguments) {
    & $command @arguments
    if ($LASTEXITCODE -ne 0) { throw "$command exited with code $LASTEXITCODE" }
}

Push-Location $root
try {
    $binary = Join-Path $root 'dist/agent-service'
    if (-not $SkipBuild) {
        $npmCache = Join-Path $env:TEMP 'agent-service-npm-cache'
        $oldCi = $env:CI
        $env:CI = 'true'
        Push-Location (Join-Path $root 'flowgram-editor')
        try {
            Invoke-Checked 'npx' @('--yes', '--cache', $npmCache, 'pnpm@10.30.2', 'install', '--frozen-lockfile')
            Invoke-Checked 'npx' @('--yes', '--cache', $npmCache, 'pnpm@10.30.2', 'build')
        } finally { Pop-Location }
        Push-Location (Join-Path $root 'web')
        try {
            Invoke-Checked 'npx' @('--yes', '--cache', $npmCache, 'pnpm@10.30.2', 'install', '--frozen-lockfile')
            if ($null -eq $oldCi) { Remove-Item Env:CI -ErrorAction SilentlyContinue } else { $env:CI = $oldCi }
            Invoke-Checked 'npx' @('--yes', '--cache', $npmCache, 'pnpm@10.30.2', 'build:sdk')
            Invoke-Checked 'npx' @('--yes', '--cache', $npmCache, 'pnpm@10.30.2', 'build')
        } finally { Pop-Location }
        New-Item -ItemType Directory -Path (Split-Path $binary) -Force | Out-Null
        $oldGoos, $oldGoarch, $oldCgo = $env:GOOS, $env:GOARCH, $env:CGO_ENABLED
        try {
            $env:GOOS = 'linux'
            $env:GOARCH = 'amd64'
            $env:CGO_ENABLED = '0'
            Invoke-Checked 'go' @('build', '-trimpath', '-ldflags=-s -w', '-o', $binary, './cmd/server')
        } finally {
            if ($null -eq $oldGoos) { Remove-Item Env:GOOS -ErrorAction SilentlyContinue } else { $env:GOOS = $oldGoos }
            if ($null -eq $oldGoarch) { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue } else { $env:GOARCH = $oldGoarch }
            if ($null -eq $oldCgo) { Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue } else { $env:CGO_ENABLED = $oldCgo }
        }
    }
    if (-not (Test-Path $binary)) { throw "Missing Linux binary: $binary" }

    $release = Get-Date -Format 'yyyyMMddHHmmssfff'
    $releaseDir = Join-Path $deployDir 'releases'
    New-Item -ItemType Directory -Path $releaseDir -Force | Out-Null
    $archive = Join-Path $releaseDir "agent-service-$release.tar.gz"
    Invoke-Checked 'tar' @('-czf', $archive, '-C', (Split-Path $binary), 'agent-service')
    $hash = (Get-FileHash $archive -Algorithm SHA256).Hash.ToLowerInvariant()
    $ssh = @('-p', [string]$target.sshPort)
    $scp = @('-P', [string]$target.sshPort)
    if ($target.identityFile) {
        $keyPath = if ([System.IO.Path]::IsPathRooted($target.identityFile)) { $target.identityFile } else { Join-Path $root $target.identityFile }
        $key = (Resolve-Path $keyPath).Path
        $ssh += @('-i', $key)
        $scp += @('-i', $key)
    }
    $remote = "$($target.user)@$($target.host)"
    $remoteArchive = "/tmp/agent-service-$release.tar.gz"
    $remoteScript = "/tmp/agent-service-deploy-$release.sh"
    $script = Join-Path $deployDir 'deploy_release.sh'
    $uploadScript = Join-Path $env:TEMP "agent-service-deploy-$release.sh"
    [System.IO.File]::WriteAllText($uploadScript, ((Get-Content $script -Raw) -replace "`r`n", "`n"), [System.Text.UTF8Encoding]::new($false))
    Invoke-Checked 'scp' ($scp + @($archive, "${remote}:$remoteArchive"))
    try {
        Invoke-Checked 'scp' ($scp + @($uploadScript, "${remote}:$remoteScript"))
        $sudo = if ($target.user -eq 'root') { '' } else { 'sudo -n ' }
        $command = "${sudo}bash '$remoteScript' '$($target.remoteRoot)' '$remoteArchive' '$hash' '$release' '$($target.serviceUser)'"
        Invoke-Checked 'ssh' ($ssh + @($remote, $command))
    } finally {
        Remove-Item $uploadScript -ErrorAction SilentlyContinue
        & ssh @ssh $remote "rm -f '$remoteArchive' '$remoteScript'" | Out-Null
    }
    Write-Host "Deployed agent-service release $release. Archive: $archive"
} finally { Pop-Location }
