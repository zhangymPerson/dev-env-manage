#!/bin/bash

# 生成测试数据的Shell脚本
# 用法: ./generate_test_data.sh

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始生成测试数据...${NC}"

# 项目列表
projects=("project-a" "project-b" "project-c" "default" "mobile-app" "web-app" "backend-service" "frontend-app")

# 环境列表
environments=("dev" "test" "prod" "staging" "default" "uat" "preprod" "local")

# 模块列表
modules=("database" "api" "frontend" "backend" "default" "auth" "payment" "notification" "logging" "cache")

# 配置键前缀
key_prefixes=("app" "db" "redis" "kafka" "elasticsearch" "mongodb" "postgres" "mysql" "nginx" "apache")

# 配置值示例
values=(
    "localhost:3306"
    "127.0.0.1:6379"
    "192.168.1.100:9092"
    "https://api.example.com"
    "secret-key-12345"
    "admin@example.com"
    "true"
    "false"
    "100"
    "500"
    "1000"
    "30m"
    "1h"
    "24h"
    "DEBUG"
    "INFO"
    "WARN"
    "ERROR"
    "/var/log/app.log"
    "/tmp/cache"
    "us-east-1"
    "eu-west-1"
    "asia-pacific-1"
)

# 计数器
count=0
total=100

echo -e "${BLUE}生成 $total 条测试配置数据...${NC}"

# 清空现有数据（可选）
# echo -e "${YELLOW}警告: 这将删除所有现有配置数据${NC}"
# read -p "是否继续? (y/N): " confirm
# if [[ $confirm != "y" && $confirm != "Y" ]]; then
#     echo "操作取消"
#     exit 0
# fi

# 生成测试数据
for ((i=1; i<=$total; i++)); do
    # 随机选择项目、环境、模块
    project=${projects[$RANDOM % ${#projects[@]}]}
    env=${environments[$RANDOM % ${#environments[@]}]}
    module=${modules[$RANDOM % ${#modules[@]}]}
    
    # 生成配置键
    key_prefix=${key_prefixes[$RANDOM % ${#key_prefixes[@]}]}
    key_suffix=$((RANDOM % 1000))
    config_key="${key_prefix}.config.${key_suffix}"
    
    # 生成配置值
    value=${values[$RANDOM % ${#values[@]}]}
    config_value="${value}-${i}"
    
    # 生成别名
    config_alias="${key_prefix}_${key_suffix}"
    auto_alias=$(echo "$config_key" | sed 's/\.//g' | cut -c1-8)
    
    # 使用dem命令添加配置
    echo -e "${YELLOW}添加配置 [$i/$total]: ${NC}项目=$project, 环境=$env, 模块=$module, 键=$config_key"
    
    ./bin/dem -p "$project" -e "$env" -m "$module" add "$config_key" "$config_value"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 成功添加配置${NC}"
        ((count++))
    else
        echo -e "${RED}✗ 添加配置失败${NC}"
    fi
    
    # 每10条数据暂停一下
    if (( i % 10 == 0 )); then
        echo -e "${BLUE}已完成 $i/$total 条数据${NC}"
        sleep 1
    fi
done

echo -e "${GREEN}测试数据生成完成!${NC}"
echo -e "成功添加 ${count}/${total} 条配置数据"

# 显示统计信息
echo -e "\n${BLUE}=== 数据统计 ===${NC}"
echo -e "项目分布:"
for project in "${projects[@]}"; do
    project_count=$(./bin/dem list | grep -c "项目: $project" || true)
    echo -e "  $project: $project_count 条配置"
done

echo -e "\n环境分布:"
for env in "${environments[@]}"; do
    env_count=$(./bin/dem list | grep -c "环境: $env" || true)
    echo -e "  $env: $env_count 条配置"
done

echo -e "\n可以使用以下命令测试:"
echo -e "  ${YELLOW}./bin/dem list${NC} - 列出所有配置"
echo -e "  ${YELLOW}./bin/dem -v list${NC} - 详细列出所有配置"
echo -e "  ${YELLOW}./bin/dem -p project-a -e dev list${NC} - 按项目和环境过滤"
echo -e "  ${YELLOW}./bin/dem get app.config.123${NC} - 获取特定配置"
echo -e "  ${YELLOW}./bin/dem delete app.config.123${NC} - 删除配置(需要确认)"

echo -e "\n${GREEN}测试数据生成脚本完成!${NC}"