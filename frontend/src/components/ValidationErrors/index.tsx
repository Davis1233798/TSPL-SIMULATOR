import React from 'react';
import { ValidationError } from '../../types/api';
import './styles.css';

interface ValidationErrorsProps {
  errors: ValidationError[];
}

const ValidationErrors: React.FC<ValidationErrorsProps> = ({ errors }) => {
  if (!errors || errors.length === 0) {
    return null;
  }

  return (
    <div className="validation-errors">
      <div className="validation-errors-header">
        <h4>後端驗證錯誤</h4>
        <span className="error-count">{errors.length} 個錯誤</span>
      </div>

      <div className="validation-errors-list">
        {errors.map((error, index) => (
          <div key={index} className="validation-error-item">
            <div className="error-icon">❌</div>
            <div className="error-content">
              {error.line > 0 && (
                <div className="error-line">行 {error.line}</div>
              )}
              {error.command && (
                <div className="error-command">命令: {error.command}</div>
              )}
              <div className="error-message">{error.message}</div>
            </div>
          </div>
        ))}
      </div>

      <div className="validation-errors-footer">
        <p>💡 提示: 請修正上述錯誤後再次嘗試渲染</p>
      </div>
    </div>
  );
};

export default ValidationErrors;
