import React from 'react';
import { renderTSPL } from '../../services/tsplApi';
import { RenderData } from '../../types/tspl';
import { ValidationError } from '../../types/api';
import { useTranslation } from 'react-i18next';
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
  const { t } = useTranslation();

  const handlePreview = async () => {
    if (!tsplCode.trim()) {
      onError(t('enterTSPLCommand'));
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
        if (response.validation_errors && response.validation_errors.length > 0) {
          if (onValidationErrors) {
            onValidationErrors(response.validation_errors);
          }
          onError(response.error || t('syntaxValidationFailed'));
        } else {
          onError(response.error || t('renderFailed'));
        }
      }
    } catch (error: any) {
      onError(error.message || t('renderFailed'));
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
        {t('previewButton')}
      </button>
      <button className="btn btn-secondary" onClick={handleClear}>
        {t('clearButton')}
      </button>
    </div>
  );
};

export default ControlPanel;
