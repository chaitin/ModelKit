import { useTranslation } from 'react-i18next'

import CustomTag, { CustomTagProps } from '../CustomTag'

type Props = {
  size?: number
  showTooltip?: boolean
  showLabel?: boolean
} & Omit<CustomTagProps, 'size' | 'tooltip' | 'icon' | 'color' | 'children'>

export const CodeTag = ({ size, showTooltip, showLabel, ...restProps }: Props) => {
  const { t } = useTranslation()

  return (
    <CustomTag
      size={size}
      color="#6366f1"
      icon="Code"
      tooltip={showTooltip ? t('models.type.code_generation') : undefined}
      {...restProps}>
      {showLabel ? t('models.type.code_generation') : ''}
    </CustomTag>
  )
}