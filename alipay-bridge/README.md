# Alipay Bridge

国内支付宝中转服务。它只负责和支付宝交互，最终充值入账仍由海外主服务的 `/api/alipay/bridge/settle` 完成。

## Build

```bash
cd alipay-bridge
go build -o alipay-bridge .
```

## Linux Package

```bash
cd alipay-bridge
./package-linux.sh
```

输出文件：

```bash
dist/alipay-bridge-linux-amd64
```

复制二进制和 `.env.example` 到 Linux 服务器后：

```bash
cp .env.example .env
vim .env
chmod +x alipay-bridge-linux-amd64
./alipay-bridge-linux-amd64
```

默认监听端口为 `3001`。

在 Windows 上交叉编译 Linux 二进制：

```powershell
cd alipay-bridge
.\package-linux.ps1
```

生成 Docker 部署文件：

```powershell
.\package-linux.ps1 all
```

`dist/` 目录会得到：

```text
alipay-bridge
Dockerfile
docker-compose.yml
.env.example
```

在 Linux 服务器进入 `dist/` 后：

```bash
docker build -t alipay-bridge:latest .
docker compose up -d
```

## Configure

```bash
cp .env.example .env
```

编辑 `.env` 后启动程序即可。程序会自动读取当前工作目录或二进制同目录下的 `.env`；如果同时设置了系统环境变量，系统环境变量优先。

也可以显式指定配置文件：

```bash
ALIPAY_BRIDGE_ENV_FILE=/etc/alipay-bridge.env ./alipay-bridge
```

## Required Env

```bash
ALIPAY_APP_ID=...
ALIPAY_PRIVATE_KEY=...
ALIPAY_PUBLIC_KEY=...
ALIPAY_BRIDGE_SECRET=...
ALIPAY_BRIDGE_OVERSEAS_SETTLE_URL=https://main.example.com/api/alipay/bridge/settle
```

## Optional Env

```bash
ALIPAY_BRIDGE_LISTEN_ADDR=:8088
ALIPAY_SANDBOX=false
ALIPAY_BRIDGE_PUBLIC_BASE_URL=https://pay-cn.example.com
ALIPAY_BRIDGE_RETURN_SUCCESS_URL=https://main.example.com/console/topup?pay=success
ALIPAY_BRIDGE_RETURN_FAIL_URL=https://main.example.com/console/topup?pay=fail
```

海外主服务需要配置：

```bash
ALIPAY_BRIDGE_ENABLED=true
ALIPAY_BRIDGE_CREATE_URL=https://pay-cn.example.com/api/alipay/create
ALIPAY_BRIDGE_SECRET=...
```
