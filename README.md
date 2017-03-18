# oss

Aliyun OSS (Object Storage Service) resource access URL generator.

阿里云 OSS 对象存储服务资源访问 URL 生成器。用于实现第三方授权访问。

## Usage

```go
import "github.com/ggicci/oss"

var (
  endpoint = "http://oss-cn-hangzhou.aliyuncs.com"
  key = "<key>"
  secret = "<secret>"
  bucket = "<bucket>"
  object = "<object>"
)

c := oss.NewClient(endpoint, key, secret)
b := c.NewBucket(bucket)

ticket := b.NewTicket("PUT", object)
ticket.Header.Set("Content-Type", "text/plain")
// let oss server check file content md5
ticket.Header.Set("Content-MD5", "eB5eJF1ptWaXm4bijSPyxw==")

if err := ticket.Sign(); err != nil {
  // ...
}
```

The ticket (in JSON format) looks like:

```json
{
    "verb": "PUT",
    "bucket": "<bucket>",
    "object": "<object>",
    "expires_at": "2017-03-18T14:15:32.248847808+08:00",
    "header": {
        "Content-Type": [
            "text/plain"
        ],
        "Content-MD5": [
            "eB5eJF1ptWaXm4bijSPyxw=="
        ],
    },
    "query": {
        "Expires": [
            "1489817732"
        ],
        "OSSAccessKeyId": [
            "<key>"
        ],
        "Signature": [
            "K63uyzyjvCJdF8N+3XZ9JQ5QMt8="
        ]
    },
    "url": "http://<bucket>.oss-cn-hangzhou.aliyuncs.com/<object>?Expires=1489817732&OSSAccessKeyId=<key>&Signature=K63uyzyjvCJdF8N%2B3XZ9JQ5QMt8%3D"
}
```
