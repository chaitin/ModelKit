import React, { useState, useMemo } from 'react';
import { Box, Chip, Stack } from '@mui/material';
import { useTheme, alpha as addOpacityToColor } from '@mui/material';
import {
  isVisionModel,
  isWebSearchModel,
  isReasoningModel,
  isFunctionCallingModel,
  isCodeModel,
  isEmbeddingModel,
  isRerankModel
} from '../utils/model';

interface ModelTagFilterProps {
  models: Array<{
    model: string;
    provider: string;
    [key: string]: any;
  }>;
  onFilteredModelsChange: (filteredModels: any[]) => void;
  selectedTags?: string[];
  onSelectedTagsChange?: (tags: string[]) => void;
}

const ModelTagFilter: React.FC<ModelTagFilterProps> = ({
  models,
  onFilteredModelsChange,
  selectedTags = [],
  onSelectedTagsChange
}) => {
  const theme = useTheme();
  const [internalSelectedTags, setInternalSelectedTags] = useState<string[]>([]);
  
  // 使用外部传入的selectedTags或内部状态
  const currentSelectedTags = selectedTags.length > 0 ? selectedTags : internalSelectedTags;
  
  // 定义可用的标签类型
  const availableTags = [
    { key: 'all', label: '全部', color: 'default' as const },
    { key: 'code', label: '代码生成', color: 'warning' as const },
    { key: 'reasoning', label: '深度思考', color: 'primary' as const },
    { key: 'function_calling', label: '工具调用', color: 'success' as const },
    { key: 'vision', label: '视觉', color: 'secondary' as const },
    { key: 'embedding', label: '向量', color: 'error' as const },
    { key: 'rerank', label: '重排', color: 'default' as const }
  ];

  // 根据模型计算每个标签的可用性
  const tagStats = useMemo(() => {
    const stats: Record<string, { available: boolean }> = {};
    
    availableTags.forEach(tag => {
      if (tag.key === 'all') {
        stats[tag.key] = { available: true };
        return;
      }
      
      const count = models.filter(model => {
        switch (tag.key) {
          case 'reasoning':
            return isReasoningModel(model.model, model.provider);
          case 'vision':
            return isVisionModel(model.model, model.provider);
          case 'websearch':
            return isWebSearchModel(model.model, model.provider);
          case 'function_calling':
            return isFunctionCallingModel(model.model, model.provider);
          case 'code':
            return isCodeModel(model.model, model.provider);
          case 'embedding':
            return isEmbeddingModel(model.model, model.provider);
          case 'rerank':
            return isRerankModel(model.model);
          default:
            return false;
        }
      }).length;
      
      stats[tag.key] = { available: count > 0 };
    });
    
    return stats;
  }, [models]);

  // 根据选中的标签筛选模型
  const filteredModels = useMemo(() => {
    if (currentSelectedTags.length === 0 || currentSelectedTags.includes('all')) {
      return models;
    }
    
    return models.filter(model => {
      return currentSelectedTags.some(tag => {
        switch (tag) {
          case 'reasoning':
            return isReasoningModel(model.model, model.provider);
          case 'vision':
            return isVisionModel(model.model, model.provider);
          case 'websearch':
            return isWebSearchModel(model.model, model.provider);
          case 'function_calling':
            return isFunctionCallingModel(model.model, model.provider);
          case 'code':
            return isCodeModel(model.model, model.provider);
          case 'embedding':
            return isEmbeddingModel(model.model, model.provider);
          case 'rerank':
            return isRerankModel(model.model);
          default:
            return false;
        }
      });
    });
  }, [models, currentSelectedTags]);

  // 当筛选结果变化时通知父组件
  React.useEffect(() => {
    onFilteredModelsChange(filteredModels);
  }, [filteredModels, onFilteredModelsChange]);

  const handleTagClick = (tagKey: string) => {
    let newSelectedTags: string[];
    
    if (tagKey === 'all') {
      newSelectedTags = [];
    } else {
      // 单选逻辑：如果点击的是已选中的tag，则取消选择；否则只选择这个tag
      if (currentSelectedTags.includes(tagKey)) {
        newSelectedTags = [];
      } else {
        newSelectedTags = [tagKey];
      }
    }
    
    if (onSelectedTagsChange) {
      onSelectedTagsChange(newSelectedTags);
    } else {
      setInternalSelectedTags(newSelectedTags);
    }
  };

  const isTagSelected = (tagKey: string) => {
    if (tagKey === 'all') {
      return currentSelectedTags.length === 0 || currentSelectedTags.includes('all');
    }
    return currentSelectedTags.includes(tagKey);
  };

  return (
    <Box>
      <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
        {availableTags.map(tag => {
          const stats = tagStats[tag.key];
          const selected = isTagSelected(tag.key);
          
          return (
            <Chip
              key={tag.key}
              label={tag.label}
              variant={selected ? 'filled' : 'outlined'}
              color={selected ? tag.color : 'default'}
              size="small"
              disabled={!stats.available}
              onClick={(e) => {
                e.stopPropagation();
                handleTagClick(tag.key);
              }}
              sx={{
                mb: 0.5,
                cursor: stats.available ? 'pointer' : 'not-allowed',
                '&:hover': stats.available ? {
                  backgroundColor: selected 
                    ? addOpacityToColor(theme.palette[tag.color === 'default' ? 'primary' : tag.color].main, 0.8)
                    : addOpacityToColor(theme.palette[tag.color === 'default' ? 'primary' : tag.color].main, 0.1)
                } : {},
                ...(selected && {
                  backgroundColor: theme.palette[tag.color === 'default' ? 'primary' : tag.color].main,
                  color: theme.palette[tag.color === 'default' ? 'primary' : tag.color].contrastText
                })
              }}
            />
          );
        })}
      </Stack>
    </Box>
  );
};

export default ModelTagFilter;