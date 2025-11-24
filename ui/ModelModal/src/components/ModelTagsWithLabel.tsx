import {
  isCodeModel,
  isEmbeddingModel,
  isFunctionCallingModel,
  isReasoningModel,
  isRerankModel,
  isVisionModel,
  isWebSearchModel,
  isAnalysisModel,
} from '../utils/model';
import { FC, memo, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { styled } from '@mui/material';
import { useTranslation } from 'react-i18next';
import {
  CodeTag,
  EmbeddingTag,
  ReasoningTag,
  RerankerTag,
  ToolsCallingTag,
  VisionTag,
  WebSearchTag,
  AnalysisTag,
} from './Tags/ModelCapabilities';

interface ModelTagsProps {
  model_id: string;
  provider: string;
  showFree?: boolean;
  showReasoning?: boolean;
  showToolsCalling?: boolean;
  size?: number;
  showLabel?: boolean;
  showTooltip?: boolean;
  style?: React.CSSProperties;
}

const ModelTagsWithLabel: FC<ModelTagsProps> = ({
  model_id,
  provider,
  showFree = true,
  showReasoning = true,
  showToolsCalling = true,
  size = 12,
  showLabel = true,
  showTooltip = true,
  style,
}) => {
  const { t } = useTranslation();
  const [shouldShowLabel, setShouldShowLabel] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const resizeObserver = useRef<ResizeObserver | null>(null);

  const maxWidth = useMemo(() => 300, []);

  useLayoutEffect(() => {
    const currentElement = containerRef.current;
    if (!showLabel || !currentElement) return;

    setShouldShowLabel(currentElement.offsetWidth >= maxWidth);

    if (currentElement) {
      resizeObserver.current = new ResizeObserver((entries) => {
        for (const entry of entries) {
          const { width } = entry.contentRect;
          setShouldShowLabel(width >= maxWidth);
        }
      });
      resizeObserver.current.observe(currentElement);
    }
    return () => {
      if (resizeObserver.current && currentElement) {
        resizeObserver.current.unobserve(currentElement);
        resizeObserver.current.disconnect();
        resizeObserver.current = null;
      }
    };
  }, [maxWidth, showLabel]);

  return (
    <Container ref={containerRef} style={style}>
      {isVisionModel(model_id, provider) && (
        <VisionTag
          size={size}
          showTooltip={showTooltip}
          showLabel={shouldShowLabel}
        />
      )}
      {isWebSearchModel(model_id, provider) && (
        <WebSearchTag
          size={size}
          showTooltip={showTooltip}
          showLabel={shouldShowLabel}
        />
      )}
      {showReasoning && isReasoningModel(model_id, provider) && (
        <ReasoningTag
          size={size}
          showTooltip={showTooltip}
          showLabel={shouldShowLabel}
        />
      )}
      {showToolsCalling && isFunctionCallingModel(model_id, provider) && (
        <ToolsCallingTag
          size={size}
          showTooltip={showTooltip}
          showLabel={shouldShowLabel}
        />
      )}
      {isCodeModel(model_id, provider) && (
        <CodeTag
          size={size}
          showTooltip={showTooltip}
          showLabel={shouldShowLabel}
        />
      )}
      {isAnalysisModel(model_id, provider) && (
        <AnalysisTag
          size={size}
          showTooltip={showTooltip}
          showLabel={shouldShowLabel}
        />
      )}
      {isEmbeddingModel(model_id, provider) && <EmbeddingTag size={size} />}
      {isRerankModel(model_id) && <RerankerTag size={size} />}
    </Container>
  );
};

const Container = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '4px',
  flexWrap: 'nowrap',
  overflowS: 'scroll',
  '&::-webkit-scrollbar': {
    display: 'none',
  },
}));

export default memo(ModelTagsWithLabel);
