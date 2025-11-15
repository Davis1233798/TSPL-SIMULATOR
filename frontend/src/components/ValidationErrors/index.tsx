import React from 'react';
import { ValidationError } from '../../types/api';
import { useTranslation } from 'react-i18next';
import './styles.css';

interface ValidationErrorsProps {
  errors: ValidationError[];
}

const ValidationErrors: React.FC<ValidationErrorsProps> = ({ errors }) => {
  const { t } = useTranslation();

  if (!errors || errors.length === 0) {
    return null;
  }

  return (
    <div className="validation-errors">
      <div className="validation-errors-header">
        <h4>{t('backendValidationErrors')}</h4>
        <span className="error-count">{errors.length} {t('errorsCount')}</span>
      </div>

      <div className="validation-errors-list">
        {errors.map((error, index) => (
          <div key={index} className="validation-error-item">
            <div className="error-icon">❌</div>
            <div className="error-content">
              {error.line > 0 && (
                <div className="error-line">{t('line')} {error.line}</div>
              )}
              {error.command && (
                <div className="error-command">{t('command')}: {error.command}</div>
              )}
              <div className="error-message">{error.message}</div>
            </div>
          </div>
        ))}
      </div>

      <div className="validation-errors-footer">
        <p>{t('hint')}</p>
      </div>
    </div>
  );
};

export default ValidationErrors;
