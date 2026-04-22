#!/bin/bash

# JWT密钥设置脚本

echo "=== JWT密钥安全设置 ==="

# 生成新的JWT密钥
echo "正在生成安全的JWT密钥..."
NEW_SECRET=$(go run config/jwt_generator.go 2>/dev/null || echo "Failed to generate")

if [ -n "$NEW_SECRET" ]; then
    echo "生成的新密钥: $NEW_SECRET"
    echo "警告: 请妥善保存此密钥，丢失后将无法验证现有令牌！"

    # 询问用户是否要使用新密钥
    read -p "是否要使用此新密钥？(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # 更新配置文件
        sed -i "s/secretKey: .*/secretKey: $NEW_SECRET/" config/settings.yml
        echo "配置文件已更新"
    else
        echo "保留当前密钥"
    fi
else
    echo "生成密钥失败，请检查配置"
fi

# 检查当前密钥的安全性
echo -e "\n=== 安全检查 ==="
CURRENT_SECRET=$(grep "secretKey:" config/settings.yml | cut -d' ' -f2)

if grep -q "123456" config/settings.yml; then
    echo "❌ 检测到弱JWT密钥: 123456"
    echo "   这可能导致安全风险，请立即更换！"
else
    if [ ${#CURRENT_SECRET} -ge 32 ]; then
        echo "✅ JWT密钥长度足够: ${#CURRENT_SECRET} 字符"
    else
        echo "❌ JWT密钥太短: ${#CURRENT_SECRET} 字符（建议至少32字符）"
    fi
fi

echo -e "\n完成！"