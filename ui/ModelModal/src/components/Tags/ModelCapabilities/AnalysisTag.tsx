import { useTranslation } from 'react-i18next'

import CustomTag, { CustomTagProps } from '../CustomTag'

type Props = {
  size?: number
  showTooltip?: boolean
  showLabel?: boolean
} & Omit<CustomTagProps, 'size' | 'tooltip' | 'icon' | 'color' | 'children'>

export const AnalysisTag = ({ size, showTooltip, showLabel, ...restProps }: Props) => {
  const { t } = useTranslation()
  return (
    <CustomTag
      size={size}
      color="#2db7f5"
      icon="分析"
      tooltip={showTooltip ? '分析' : undefined}
      {...restProps}>
      {showLabel ? '分析' : ''}
    </CustomTag>
  )
}
