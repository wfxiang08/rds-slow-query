# Rds Slow Query分析&邮件发送

## 编译
```
git clone https://github.com/wfxiang08/rds_slowq_query.git
cd rds_slowq_query
source start_env.sh
glide install
go build cmds/tool_slow_query_analyze.go
```

## 使用
```
./tool_slow_query_analyze -conf conf/db.toml
```
