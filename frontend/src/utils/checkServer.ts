import { api } from '@/utils/api';


export const checkServerHealth = async (): Promise<boolean> => {
    try {
      const response = await api.getVersion();
      return response.status === 200;
    } catch (error) {
      return false;
    }
  };
  