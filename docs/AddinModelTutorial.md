# 接入本地部署大模型教程

# 安装部署平台(Ollama/Xinference/GPUStack)

请按照您想要接入的平台的教程部署, 注意 pandawiki只支持linux系统部署

[ollama安装教程](https://docs.ollama.com/quickstart)

[xinference安装教程](https://inference.readthedocs.io/zh-cn/v1.2.0/getting_started/installation.html)

[gpustack安装教程](https://docs.gpustack.ai/latest/quickstart/)

# 确认大模型平台安装成功

## 将大模型平台的监听IP设置为0.0.0.0

### ollama

1.  通过执行 `systemctl edit ollama.service`编辑 systemd 服务文件
    
2.  在 `[Service]`部分下，添加一行 `Environment`
    

```ini
[Service]
Environment="OLLAMA_HOST=0.0.0.0:11434"
```

1.  保存并退出编辑器。
    
2.  重新加载 systemd 并重启 Ollama：
    

```bash
systemctl daemon-reload
systemctl restart ollama
```

### xinference

在启动 Xinference 时加上 `-H 0.0.0.0` 参数:

```plaintext
xinference-local -H 0.0.0.0
```

## 获取部署机器的ip: 在命令行中输入 ip addr, 通常为 eth0或wlan0 网卡中的inet 后的ip

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/mxPOG5zpEKkaJnKa/img/483323b8-d6e7-4f30-ab82-e500cfb170c3.png)

## 获取模型列表

请替换命令中的端口号为 ollama的默认端口号: 11434, xinference的默认端口号: 9997, gpustack的默认端口号: 80

```c++
curl -X GET \
  http://部署机器的ip:端口/v1/models \
  -H "Content-Type: application/json"
  
```
示例响应
```c++

{
  "object": "list",
  "data": [
    {
      "id": "gemma3:latest",
      "object": "model",
      "created": 1755516177,
      "owned_by": "library"
    }
  ]
}
```

## 检查模型是否可以使用
对话模型
```c++
  curl -X POST \
  http://部署机器的ip:端口/v1/chat/completions \
  -H "Authorization: Bearer 您设置的API Key,没有可以去掉这行" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "模型列表中您想要配置的模型id",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```
向量模型
```c++
  curl -X POST \
  http://部署机器的ip:端口/v1/embeddings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 您设置的API Key,没有可以去掉这行" \
  -d '{
    "model": "模型列表中您想要配置的模型id",
    "input": "hello, nice to meet you , and you?",
    "encoding_format": "float"
  }'
```
重排序模型
```c++
  curl -X POST \
  http://部署机器的ip:端口/v1/rerank \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 您设置的API Key,没有可以去掉这行" \
  -d '{
    "model": "模型列表中您想要配置的模型id",
    "documents": ["hello"],
	"query": "test"
  }'
```

# 配置模型

## 选择供应商

**对话模型**

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/23a5b54d-be99-4cd3-81b7-b1e47ebca6c2.png)

**向量/重排序模型**

注意 ollama不支持重排序模型!

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/410271ce-dfa0-4ec5-bc89-9a0c2328099d.png)

## 输入API地址与API Key![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/bb3327ac-deda-4b6d-bf6c-970f582b469a.png)

1.  API地址为`http://curl中的ip:curl中的端口`
    
2.  选择其它供应商时, 还需要输入之前curl中您输入的模型名称
    

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/9284fe77-8d1d-40d6-8689-06adb5c82767.png)

## 选择模型

注意向量/重排只能选择对应标签下的模型

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/39e3c799-5617-4b6d-bd3b-089dd47590f1.png)

## 确认保存

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/c3ae5021-b4a8-4228-b03f-bfe48f0e86cb.png)

配置成功后会弹出“修改成功”的提示

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/54Lq35ojy3gLXl7E/img/a31a2600-bea0-478c-8a48-c473013035a4.png)

# 按照上述流程执行, 依然配置失败怎么办?

请附上报错的截图提Issue, 开发者会及时解答您的问题