import { useContext } from 'react';
import { AuthContext, CommonContext } from '@/context';
import { useModelContext } from '@/providers/ModelProvider';

export const useAuthContext = () => {
  return useContext(AuthContext);
};

export const useCommonContext = () => {
  // 使用ModelProvider提供的数据
  return useModelContext();
};