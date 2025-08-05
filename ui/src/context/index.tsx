import { createContext } from 'react';
import { DomainModel } from '@/api/types';

// 暂时定义缺失的类型
interface DomainUser {
  id?: string;
  username?: string;
}

interface DomainAdminUser extends DomainUser {
  role?: string;
}

export const AuthContext = createContext<
  [
    DomainUser | DomainAdminUser | null,
    {
      loading: boolean;
      setUser: (user: DomainUser | DomainAdminUser) => void;
      refreshUser: () => void;
    }
  ]
>([
  null,
  {
    setUser: () => {},
    loading: true,
    refreshUser: () => {},
  },
]);

export const CommonContext = createContext<{
  coderModel: DomainModel[];
  llmModel: DomainModel[];
  isConfigModel: boolean;
  modelLoading: boolean;
  refreshModel: () => void;
}>({
  coderModel: [],
  llmModel: [],
  isConfigModel: false,
  modelLoading: false,
  refreshModel: () => {},
});