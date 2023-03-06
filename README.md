# streak

## 数据模式
监听所有上行流量，并进行统计。

* dns 模式，根据所有 dns 明文报文，解析出 (sourceIp,hostName) 对，同时记录下解析出的 (ip,hostName)对用于后续兜底策略。
* http 模式，根据所有 http 明文报文，解析出所有带域名的 (sourceIp,hostName)对。
* tls 模式，根据所有 tls 报文，根据其 sni（tls1.3后使用 esni 将无法获取），解析出所有带域名的（sourceIp，sni）对。
* transport 模式，根据所有的 transport 层（tcp+udp）报文，结合 dns 模式中顺带记录的 （ip，hostName）对，反向推测出当前连接可能的（sourceIp，hostName）对。

## 网络模式
* 本机模式
* 主路由模式
* 旁路由模式

## 域名展示模式
* 四层模式下支持猜测域名的缩略展示。
* 四层模式以上支持真实域名的完整展示和缩略展示。
* 缩略展示支持

#