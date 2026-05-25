# Alipay Bridge Windows Build Helper
# Usage: .\package-linux.ps1 [build|docker|all]

param(
  [Parameter(Position = 0)]
  [ValidateSet('build', 'docker', 'all')]
  [string]$Action = 'build'
)

$ErrorActionPreference = 'Stop'

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DistDir = Join-Path $RootDir 'dist'
$Output = Join-Path $DistDir 'alipay-bridge'
$GoCache = Join-Path $DistDir '.gocache'

function Write-Bridge-Header($msg) {
  Write-Host ''
  Write-Host '+------------------------------------------------------------+' -ForegroundColor Cyan
  Write-Host ('| ' + $msg.PadRight(58) + ' |') -ForegroundColor Cyan
  Write-Host '+------------------------------------------------------------+' -ForegroundColor Cyan
}

function Invoke-Checked($Command, [string[]]$CommandArgs) {
  & $Command @CommandArgs
  if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
  }
}

function Build-LinuxBinary {
  Write-Bridge-Header 'BUILDING ALIPAY BRIDGE: LINUX AMD64'
  Set-Location $RootDir

  Write-Host '>> Preparing dist directory...'
  New-Item -ItemType Directory -Force -Path $DistDir | Out-Null
  New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
  $LegacyOutput = Join-Path $DistDir 'alipay-bridge-linux-amd64'
  if (Test-Path $LegacyOutput) {
    Remove-Item -LiteralPath $LegacyOutput -Force
  }

  $env:GOCACHE = $GoCache

  Write-Host '>> Running tests...'
  Invoke-Checked 'go' @('test', './...')

  Write-Host '>> Compiling linux/amd64 binary from Windows...'
  $env:CGO_ENABLED = '0'
  $env:GOOS = 'linux'
  $env:GOARCH = 'amd64'
  Invoke-Checked 'go' @('build', '-buildvcs=false', '-trimpath', '-ldflags=-s -w', '-o', $Output, '.')

  Write-Host ''
  Write-Host "[SUCCESS] Binary created: $Output" -ForegroundColor Green
  Write-Host '>> Default listen address is configured as :3001' -ForegroundColor Yellow
}

function Write-DockerArtifacts {
  Write-Bridge-Header 'WRITING DOCKER DEPLOYMENT FILES'
  Set-Location $RootDir
  New-Item -ItemType Directory -Force -Path $DistDir | Out-Null

  $dockerfile = @'
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget \
    && update-ca-certificates

COPY alipay-bridge /alipay-bridge
RUN chmod +x /alipay-bridge

EXPOSE 3001
WORKDIR /data

ENTRYPOINT ["/alipay-bridge"]
'@

  $compose = @'
version: '3.4'

services:
  alipay-bridge:
    image: alipay-bridge:latest
    container_name: alipay-bridge
    restart: always
    ports:
      - "3001:3001"
    environment:
      - ALIPAY_BRIDGE_LISTEN_ADDR=:3001
      - ALIPAY_APP_ID=
      - ALIPAY_PRIVATE_KEY=
      - ALIPAY_PUBLIC_KEY=
      - ALIPAY_SANDBOX=false
      - ALIPAY_BRIDGE_PUBLIC_BASE_URL=https://pay-cn.example.com
      - ALIPAY_BRIDGE_OVERSEAS_SETTLE_URL=https://main.example.com/api/alipay/bridge/settle
      - ALIPAY_BRIDGE_SECRET=
      - ALIPAY_BRIDGE_RETURN_SUCCESS_URL=https://main.example.com/console/topup?pay=success
      - ALIPAY_BRIDGE_RETURN_FAIL_URL=https://main.example.com/console/topup?pay=fail
      - TZ=Asia/Shanghai
    networks:
      - alipay-bridge-network
    healthcheck:
      test: ["CMD-SHELL", "wget -q -O - http://localhost:3001/health | grep -o 'ok' || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  alipay-bridge-network:
    driver: bridge
'@

  Set-Content -Path (Join-Path $DistDir 'Dockerfile') -Value $dockerfile -Encoding UTF8
  Set-Content -Path (Join-Path $DistDir 'docker-compose.yml') -Value $compose -Encoding UTF8
  Copy-Item -Path (Join-Path $RootDir '.env.example') -Destination (Join-Path $DistDir '.env.example') -Force

  Write-Host "[SUCCESS] Docker files created in: $DistDir" -ForegroundColor Green
  Write-Host '>> Files: alipay-bridge, Dockerfile, docker-compose.yml, .env.example' -ForegroundColor Yellow
}

switch ($Action) {
  'build' {
    Build-LinuxBinary
  }
  'docker' {
    Write-DockerArtifacts
  }
  'all' {
    Build-LinuxBinary
    Write-DockerArtifacts
  }
}
