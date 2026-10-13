# 燧原容器工具包 (enflame-container-toolkit)

> **语言**: [English](README.md) | 中文

## 目录

- [概述](#概述)
- [项目结构](#项目结构)
- [安装](#安装)
  - [前置条件](#前置条件)
  - [从源码构建](#从源码构建)
- [快速开始](#快速开始)
  - [构建二进制](#构建二进制)
  - [配置容器运行时](#配置容器运行时)
  - [运行 GCU 容器](#运行-gcu-容器)
- [核心特性](#核心特性)
- [版本兼容性](#版本兼容性)
- [硬件支持](#硬件支持)
- [故障排查](#故障排查)
- [许可证](#许可证)

## 概述

**燧原容器工具包 (Enflame Container Toolkit)** 使用户能够构建并运行 GCU
加速容器。它提供了一个 OCI 兼容的容器运行时、一个运行时钩子、一个命令行工具，
以及 CDI（容器设备接口，Container Device Interface）支持，可自动配置容器以访问
燧原 GCU（通用计算单元，General Compute Unit）设备。它是 `nvidia-container-toolkit`
在燧原上的对应实现，并与底层的 `libenflame-container` 库协同工作。

该工具包提供：

- **OCI 兼容运行时** —— `enflame-container-runtime` 封装底层 OCI 运行时，在容器启动时将 GCU 设备注入容器。
- **运行时钩子与 CDI 支持** —— 运行时钩子（`enflame-container-runtime-hook`）与容器设备接口路径（`enflame-container-runtime.cdi`）自动为容器准备 GCU 访问。
- **配置简单** —— `enflame-ctk` 命令行工具用于配置容器运行时，使容器可以透明地获取 GCU 设备。
- **兼容常见运行时** —— 可与 containerd 和 Docker 集成，并与 `libenflame-container` 库及 `enflame-docker` 使用层互操作。

## 项目结构

```
enflame-container-toolkit/
├── cmd/                        # 各二进制的入口
│   ├── enflame-container-runtime/        # OCI 兼容容器运行时
│   ├── enflame-container-runtime-hook/   # 运行时钩子（启动前配置）
│   ├── enflame-container-runtime.cdi/    # 基于 CDI 的运行时路径
│   ├── enflame-container-runtime.legacy/ # 传统运行时路径
│   └── enflame-ctk/                      # 命令行配置工具
├── internal/                   # 私有应用包
├── pkg/                        # 可复用的库包
├── deployments/                # 部署清单与资源
├── docker/                     # Docker 构建资源
├── packaging/                  # 操作系统打包文件
├── scripts/                    # 辅助脚本
├── test/                       # 测试
├── tools/                      # 开发工具
├── go.mod                      # Go 模块定义
├── Makefile                    # 构建全部二进制
└── versions.mk                 # 版本定义
```

## 安装

### 前置条件

- 一台已安装燧原 GCU 驱动的燧原 GCU 设备。
- 一个容器运行时 —— containerd 或 Docker。
- 用于从源码构建的 [Go](https://go.dev/) 及可用的工具链。
- 主机上可用的底层 `libenflame-container` 库。

### 从源码构建

```bash
# 构建运行时、钩子、CDI、传统运行时和 CLI 二进制
make

# 清理构建产物
make clean
```

构建会产出 `cmd/` 下定义的运行时、钩子和 CLI 二进制：
`enflame-container-runtime`、`enflame-container-runtime-hook`、
`enflame-container-runtime.cdi`、`enflame-container-runtime.legacy` 和
`enflame-ctk`。

## 快速开始

### 构建二进制

```bash
make
```

### 配置容器运行时

使用 `enflame-ctk` 命令行工具将燧原运行时注册到你的容器运行时（containerd 或
Docker）。配置完成后，容器运行时会将 GCU 工作负载路由到
`enflame-container-runtime`，由其通过运行时钩子 / CDI 路径注入 GCU 设备。

```bash
# 使用燧原运行时配置容器运行时（示意）
enflame-ctk runtime configure

# 重启容器运行时使改动生效
# （根据你的环境，例如 containerd 或 docker）
```

### 运行 GCU 容器

运行时配置完成后，使用燧原运行时启动容器。以这种方式启动的容器会自动注入 GCU
设备，可在容器内部访问燧原 GCU。

```bash
# 使用燧原运行时启动容器（示意）
docker run --runtime=enflame --rm -it <your-image> <command>
```

## 核心特性

| 特性 | 说明 |
| --- | --- |
| OCI 运行时 | `enflame-container-runtime` 是 OCI 兼容运行时，向容器注入 GCU 设备 |
| 运行时钩子 | `enflame-container-runtime-hook` 在容器启动时为其准备 GCU 访问 |
| CDI 支持 | 容器设备接口路径（`enflame-container-runtime.cdi`）实现声明式设备注入 |
| enflame-ctk CLI | 配置容器运行时与容器，使其透明访问 GCU |
| 运行时集成 | 兼容 containerd 与 Docker |
| 库互操作 | 基于底层 `libenflame-container` 库与 `enflame-docker` 使用层构建 |

## 版本兼容性

燧原容器工具包配置对燧原 GCU 驱动所暴露的 GCU 设备的访问，因此其行为随所安装的
驱动而定。请始终使工具包版本与主机上部署的 GCU 驱动匹配，并确保
`libenflame-container` 库与该驱动兼容。

## 硬件支持

支持提供燧原 GCU 驱动的燧原 GCU 设备。运行时、钩子和 CLI 可跨 GCU 产品代际工作；
容器内可用的设备取决于主机上安装的驱动。

## 故障排查

- **容器看不到 GCU 设备** —— 确认主机已加载燧原 GCU 驱动，且已通过 `enflame-ctk`
  将燧原运行时注册到容器运行时；配置后重启 containerd/Docker。
- **找不到 `enflame` 运行时** —— 确认运行时二进制已安装并位于预期路径，且容器运行时
  配置正确引用了 `enflame-container-runtime`。
- **容器启动时报库错误** —— 确认 `libenflame-container` 库已安装且与运行中的 GCU
  驱动兼容。

## 许可证

燧原容器工具包基于 **Apache License 2.0** 许可证发布。版权所有 (c) Enflame
Technology。保留所有权利。详见 [LICENSE](LICENSE)。
