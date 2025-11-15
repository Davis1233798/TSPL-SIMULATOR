import React from 'react';
import { renderTSPL } from '../../services/tsplApi';
import { RenderData } from '../../types/tspl';
import { ValidationError } from '../../types/api';
import './styles.css';

interface ControlPanelProps {
  tsplCode: string;
  onRenderDataUpdate: (data: RenderData | null) => void;
  onLoadingChange: (loading: boolean) => void;
  onError: (error: string) => void;
  onValidationErrors?: (errors: ValidationError[]) => void;
}

const ControlPanel: React.FC<ControlPanelProps> = ({
  tsplCode,
  onRenderDataUpdate,
  onLoadingChange,
  onError,
  onValidationErrors,
}) => {
  const handlePreview = async () => {
    if (!tsplCode.trim()) {
      onError('請輸入 TSPL 命令');
      return;
    }

    try {
      onLoadingChange(true);
      onError('');
      if (onValidationErrors) onValidationErrors([]);

      const response = await renderTSPL(tsplCode);

      if (response.success && response.data) {
        onRenderDataUpdate(response.data);
      } else {
        // 處理驗證錯誤
        if (response.validation_errors && response.validation_errors.length > 0) {
          if (onValidationErrors) {
            onValidationErrors(response.validation_errors);
          }
          onError(response.error || 'TSPL 語法驗證失敗');
        } else {
          onError(response.error || '渲染失敗');
        }
      }
    } catch (error: any) {
      onError(error.message || '渲染失敗');
      onRenderDataUpdate(null);
    } finally {
      onLoadingChange(false);
    }
  };

  const handleClear = () => {
    onRenderDataUpdate(null);
    onError('');
    if (onValidationErrors) onValidationErrors([]);
  };

  return (
    <div className="control-panel">
      <button className="btn btn-primary" onClick={handlePreview}>
        預覽
      </button>
      <button className="btn btn-secondary" onClick={handleClear}>
        清除
      </button>
    </div>
  );
};

export default ControlPanel;
