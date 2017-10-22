# auserver

> A API server for mip log data.

## API

### /list/[类别]
> 各类别及含义


# auserver

> A API server for mip log data.

## API

### /list/[类别]
> 类别及含义

- domain
    - [具体域名]

- tags
    - [具体组件名]
    - all
    - core
    - official
    - plat
    - unuse

a.g：`/list/tags/core` 则会列出所有核心组件名称

### /api/[类别]
> API接口，返回json数据

类别及含义
- tags
- tagsinfo
- count
- domains
- select
- tagsbar
- barcount
- tagtotal

a.g：`/api/tags` 则会列出组件名及引用数

### /api/[类别]
> API接口，返回json数据

类别及含义
- tags
- tagsinfo
- count
- domains
- select
- tagsbar
- barcount
- tagtotal
