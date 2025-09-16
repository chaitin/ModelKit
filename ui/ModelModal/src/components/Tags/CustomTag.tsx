import { Clear } from '@mui/icons-material';
import { Tooltip } from '@mui/material';
import { CSSProperties, FC, memo, useMemo } from 'react';
import { styled } from '@mui/material';

export interface CustomTagProps {
  icon?: React.ReactNode;
  children?: React.ReactNode | string;
  color: string;
  size?: number;
  style?: CSSProperties;
  tooltip?: string;
  closable?: boolean;
  onClose?: () => void;
  onClick?: () => void;
  disabled?: boolean;
  inactive?: boolean;
}

const CustomTag: FC<CustomTagProps> = ({
  children,
  icon,
  color,
  size = 12,
  style,
  tooltip,
  closable = false,
  onClose,
  onClick,
  disabled,
  inactive,
}) => {
  const actualColor = inactive ? '#aaaaaa' : color;
  const tagContent = useMemo(
    () => (
      <Tag
        $color={actualColor}
        $size={size}
        $closable={closable}
        onClick={disabled ? undefined : onClick}
        style={{
          cursor: disabled ? 'not-allowed' : onClick ? 'pointer' : 'auto',
          ...style,
        }}
      >
        {icon && icon} {children}
        {closable && (
          <CloseIcon
            $size={size}
            $color={actualColor}
            onClick={(e) => {
              e.stopPropagation();
              onClose?.();
            }}
          />
        )}
      </Tag>
    ),
    [
      actualColor,
      children,
      closable,
      disabled,
      icon,
      onClick,
      onClose,
      size,
      style,
    ]
  );

  return tooltip ? (
    <Tooltip title={tooltip} placement='top'>
      {tagContent}
    </Tooltip>
  ) : (
    tagContent
  );
};

export default memo(CustomTag);

const Tag = styled('div')<{
  $color: string;
  $size: number;
  $closable: boolean;
}>(({ $color, $size, $closable }) => ({
  display: 'inline-flex',
  alignItems: 'center',
  gap: '4px',
  padding: `${$size / 3}px ${$size * 0.8}px`,
  paddingRight: `${$closable ? $size * 1.8 : $size * 0.8}px`,
  borderRadius: '99px',
  color: $color,
  backgroundColor: `${$color}20`,
  fontSize: $size,
  lineHeight: 1,
  whiteSpace: 'nowrap',
  position: 'relative',
  '& .iconfont': {
    fontSize: $size,
    color: $color,
  },
}));

const CloseIcon = styled(Clear)<{ $size: number; $color: string }>(
  ({ $size, $color }) => ({
    cursor: 'pointer',
    fontSize: `${$size * 0.8}px`,
    color: $color,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    position: 'absolute',
    right: `${$size * 0.2}px`,
    top: `${$size * 0.2}px`,
    bottom: `${$size * 0.2}px`,
    borderRadius: '99px',
    transition: 'all 0.2s ease',
    aspectRatio: '1 / 1',
    lineHeight: 1,
    '&:hover': {
      backgroundColor: '#da8a8a',
      color: '#ffffff',
    },
  })
);
