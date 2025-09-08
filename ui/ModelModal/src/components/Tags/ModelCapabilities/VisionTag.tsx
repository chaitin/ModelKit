import { useTranslation } from 'react-i18next'

import CustomTag, { CustomTagProps } from '../CustomTag'

type Props = {
  size?: number
  showTooltip?: boolean
  showLabel?: boolean
} & Omit<CustomTagProps, 'size' | 'tooltip' | 'icon' | 'color' | 'children'>

export const VisionTag = ({ size, showTooltip, showLabel, ...restProps }: Props) => {
  const { t } = useTranslation()

  return (
    <CustomTag
      size={size}
      color="#00b96b"
      icon="视觉"
      tooltip={showTooltip ? t('models.type.vision') : undefined}
      {...restProps}>
      {showLabel ? t('models.type.vision') : ''}
    </CustomTag>
  )
}
