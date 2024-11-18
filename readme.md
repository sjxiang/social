

# social

> 社交网站(注册、激活、登录、注销、关注、取消关注、发帖、评论、Feed 流)
> 会员订阅
> 使用 ffmpeg 将 mp4 转为 hls 格式


## 技术方案

> 接口文档, 基于 go-swagger 生成 API 文档, 只需一行注释即可生成文档

> 日志, 基于 zap 实现日志记录, 支持文件\控制台输出

> validator, 简单易用的校验器, 只需添加 tag 即可校验

> 邮件, 走 sendgrid 这种第三方代理, 或者利用 QQ 邮箱 pop3 协议

> 认证, jwt.MapClaims 或者 jwt.RegisteredClaims

> 限流, 固定滑动窗口

> 数据库, MySQL, ❌ ORM, ✅ raw SQL

> 路由, chi


# 表

<!-- 

用户

激活码

会员订阅

关注

帖子

评论

 -->