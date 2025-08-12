// 本地化消息配置
export const LOCALE_MESSAGES: Record<'zh-CN' | 'en-US', Record<string, string>> = {
  'zh-CN': {
    // 标题
    'addChatModel': '添加 Chat 模型',
    'addEmbeddingModel': '添加 Embedding 模型',
    'addRerankModel': '添加 Rerank 模型',
    'editChatModel': '修改 Chat 模型',
    'editEmbeddingModel': '修改 Embedding 模型',
    'editRerankModel': '修改 Rerank 模型',
    
    // 标签
    'modelProvider': '模型供应商',
    'apiAddress': 'API 地址',
    'apiSecret': 'API Secret',
    'apiVersion': 'API Version',
    'modelName': '模型名称',
    'header': 'Header',
    'save': '保存',
    'cancel': '取消',
    'getModelList': '获取模型列表',
    'viewDocumentation': '查看文档',
    
    // 验证消息
    'urlRequired': 'URL 不能为空',
    'secretRequired': 'API Secret 不能为空',
    'modelNameRequired': '模型名称不能为空',
    'modelNameHint': '需要与模型供应商提供的名称完全一致，不要随便填写',
    
    // 成功/错误消息
    'addSuccess': '添加成功',
    'updateSuccess': '修改成功',
    'testSuccess': '测试成功',
    'testFailed': '测试失败',
    
    // 其他
    'other': '其他',
    'required': '必需',
  },
  'en-US': {
    // Titles
    'addChatModel': 'Add Chat Model',
    'addEmbeddingModel': 'Add Embedding Model',
    'addRerankModel': 'Add Rerank Model',
    'editChatModel': 'Edit Chat Model',
    'editEmbeddingModel': 'Edit Embedding Model',
    'editRerankModel': 'Edit Rerank Model',
    
    // Labels
    'modelProvider': 'Model Provider',
    'apiAddress': 'API Address',
    'apiSecret': 'API Secret',
    'apiVersion': 'API Version',
    'modelName': 'Model Name',
    'header': 'Header',
    'save': 'Save',
    'cancel': 'Cancel',
    'getModelList': 'Get Model List',
    'viewDocumentation': 'View Documentation',
    
    // Validation messages
    'urlRequired': 'URL is required',
    'secretRequired': 'API Secret is required',
    'modelNameRequired': 'Model name is required',
    'modelNameHint': 'Must match exactly with the name provided by the model provider',
    
    // Success/Error messages
    'addSuccess': 'Added successfully',
    'updateSuccess': 'Updated successfully',
    'testSuccess': 'Test successful',
    'testFailed': 'Test failed',
    
    // Others
    'other': 'Other',
    'required': 'Required',
  },
};

// 获取本地化消息
export const getLocaleMessage = (key: string, language: 'zh-CN' | 'en-US' = 'zh-CN'): string => {
  return LOCALE_MESSAGES[language]?.[key] || key;
};

// 获取标题映射
export const getTitleMap = (language: 'zh-CN' | 'en-US' = 'zh-CN') => ({
  chat: getLocaleMessage('addChatModel', language),
  embedding: getLocaleMessage('addEmbeddingModel', language),
  rerank: getLocaleMessage('addRerankModel', language),
}); 