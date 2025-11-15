import React, { useEffect, useState } from 'react';
import { validateTSPL, ValidationResult } from '../../services/tsplValidator';
import { useTranslation } from 'react-i18next';
import './styles.css';

interface SyntaxCheckerProps {
  tsplCode: string;
}

const SyntaxChecker: React.FC<SyntaxCheckerProps> = ({ tsplCode }) => {
  const { t } = useTranslation();
  const [validation, setValidation] = useState<ValidationResult | null>(null);

  useEffect(() => {
    if (tsplCode.trim()) {
      const result = validateTSPL(tsplCode);
      setValidation(result);
    } else {
      setValidation(null);
    }
  }, [tsplCode]);

  if (!validation || (!validation.errors.length && !validation.warnings.length)) {
    return null;
  }

  return (
    <div className="syntax-checker">
      <div className="syntax-checker-header">
        <h4>{t('syntaxCheck')}</h4>
        {validation.isValid ? (
          <span className="status-badge status-success">{t('noErrors')}</span>
        ) : (
          <span className="status-badge status-error">✗ {validation.errors.length} {t('errorsCount')}</span>
        )}
      </div>

      <div className="syntax-checker-content">
        {validation.errors.map((error, index) => (
          <div key={`error-${index}`} className="validation-message validation-error">
            <span className="message-icon">❌</span>
            <div className="message-content">
              {error.line > 0 && <span className="line-number">{t('line')} {error.line}:</span>}
              <span className="message-text">{error.message}</span>
            </div>
          </div>
        ))}

        {validation.warnings.map((warning, index) => (
          <div key={`warning-${index}`} className="validation-message validation-warning">
            <span className="message-icon">⚠️</span>
            <div className="message-content">
              {warning.line > 0 && <span className="line-number">{t('line')} {warning.line}:</span>}
              <span className="message-text">{warning.message}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SyntaxChecker;
