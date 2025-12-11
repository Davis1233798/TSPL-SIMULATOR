import React from 'react';
import { useTheme } from '../ThemeContext';
import { useTranslation } from 'react-i18next';
import './Header.css';
import WebtraficAd from './WebtraficAd';

const Header: React.FC = () => {
  const { theme, toggleTheme } = useTheme();
  const { t, i18n } = useTranslation();

  const languages = [
    { code: 'zh-TW', name: '繁體中文', flag: '🇹🇼' },
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'ja', name: '日本語', flag: '🇯🇵' }
  ];

  return (
    <header className="app-header">
      <div className="header-content">
        <div className="header-title">
          <h1>{t('title')}</h1>
          <p>{t('subtitle')}</p>
        </div>

        <WebtraficAd />

        <div className="header-controls">
          {/* Language Selector */}
          <div className="language-selector">
            <select
              value={i18n.language}
              onChange={(e) => i18n.changeLanguage(e.target.value)}
              aria-label={t('language')}
            >
              {languages.map(lang => (
                <option key={lang.code} value={lang.code}>
                  {lang.flag} {lang.name}
                </option>
              ))}
            </select>
          </div>

          {/* Theme Toggle */}
          <button
            className="theme-toggle"
            onClick={toggleTheme}
            aria-label={theme === 'light' ? t('darkMode') : t('lightMode')}
            title={theme === 'light' ? t('darkMode') : t('lightMode')}
          >
            {theme === 'light' ? '🌙' : '☀️'}
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;
