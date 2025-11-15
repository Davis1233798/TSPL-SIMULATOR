// TSPL API 服務 - 支援後端 API 和前端 fallback
import axios from 'axios';
import { RenderResponse, ExamplesResponse } from '../types/api';
import { mockRenderTSPL, mockGetExamples } from './mockApi';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

// 檢查後端是否可用
let backendAvailable = false;

export const healthCheck = async (): Promise<boolean> => {
  try {
    const response = await axios.get(`${API_BASE_URL}/health`, { timeout: 2000 });
    backendAvailable = response.data.status === 'ok';
    return backendAvailable;
  } catch (error) {
    backendAvailable = false;
    return false;
  }
};

// 渲染 TSPL - 優先使用後端,失敗則使用前端
export const renderTSPL = async (tsplCode: string): Promise<RenderResponse> => {
  // 先嘗試使用後端
  if (backendAvailable) {
    try {
      const response = await axios.post(`${API_BASE_URL}/render`, {
        tspl_code: tsplCode,
      }, { timeout: 5000 });

      if (response.data.success) {
        return response.data;
      }
    } catch (error) {
      console.warn('後端渲染失敗,使用前端 fallback:', error);
      backendAvailable = false;
    }
  }

  // 後端不可用或失敗,使用前端實作
  return mockRenderTSPL(tsplCode);
};

// 取得範例 - 優先使用後端,失敗則使用前端
export const getExamples = async (): Promise<ExamplesResponse> => {
  // 先嘗試使用後端
  if (backendAvailable) {
    try {
      const response = await axios.get(`${API_BASE_URL}/examples`, { timeout: 5000 });

      if (response.data.success) {
        return response.data;
      }
    } catch (error) {
      console.warn('後端獲取範例失敗,使用前端 fallback:', error);
      backendAvailable = false;
    }
  }

  // 後端不可用或失敗,使用前端實作
  return mockGetExamples();
};

// 取得範例詳情
export const getExampleDetail = async (id: string): Promise<{ success: boolean; code?: string; error?: string }> => {
  if (backendAvailable) {
    try {
      const response = await axios.get(`${API_BASE_URL}/examples/${id}`, { timeout: 5000 });
      return response.data;
    } catch (error) {
      console.warn('後端獲取範例詳情失敗:', error);
    }
  }

  return { success: false, error: '範例不存在' };
};

// 發布 MQTT 訊息
export const publishMQTT = async (topic: string, message: any): Promise<{ success: boolean; error?: string }> => {
  if (!backendAvailable) {
    return { success: false, error: 'Backend not available' };
  }

  try {
    const response = await axios.post(`${API_BASE_URL}/mqtt/publish`, {
      topic,
      message,
    }, { timeout: 5000 });

    return response.data;
  } catch (error: any) {
    return {
      success: false,
      error: error.response?.data?.error || error.message
    };
  }
};

// 初始化 - 檢查後端狀態
healthCheck();
