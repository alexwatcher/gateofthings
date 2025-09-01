'use client';

import { useTranslations } from 'next-intl';
import React from 'react';

export default function LoginPage() {
  const t = useTranslations();

  return (
    <>
      <style jsx>{`
        input::placeholder {
          color: rgba(0, 255, 0, 0.5); /* pale green */
          text-shadow: none; /* no glow for placeholders */
        }
      `}</style>

      <div
        style={{
          minHeight: '100vh',
          width: '100vw',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: 'black',
          fontFamily: 'monospace',
          color: '#00ff00',
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '1.5rem',
            userSelect: 'none',
          }}
        >
          <h1
            style={{
              fontSize: '2.5rem',
              fontWeight: 'bold',
              letterSpacing: '2px',
              textShadow: '0 0 10px #00ff00, 0 0 20px #00ff00',
              margin: 0,
            }}
          >
            {t('title')}
          </h1>

          <form
            style={{
              background: 'rgba(0, 0, 0, 0.8)',
              padding: '2rem 2.5rem',
              borderRadius: '1rem',
              boxShadow: '0 0 20px #00ff00',
              display: 'flex',
              flexDirection: 'column',
              gap: '1.2rem',
              minWidth: '320px',
              border: '1px solid #00ff00',
            }}
          >
            <input
              type="text"
              placeholder={t('loginPlaceholder')}
              style={{
                padding: '0.75rem 1rem',
                borderRadius: '0.5rem',
                border: '1px solid #00ff00',
                fontSize: '1rem',
                color: '#00ff00',
                background: 'black',
                outline: 'none',
                textShadow: '0 0 5px #00ff00',
              }}
            />
            <input
              type="password"
              placeholder={t('passwordPlaceholder')}
              style={{
                padding: '0.75rem 1rem',
                borderRadius: '0.5rem',
                border: '1px solid #00ff00',
                fontSize: '1rem',
                color: '#00ff00',
                background: 'black',
                outline: 'none',
                textShadow: '0 0 5px #00ff00',
              }}
            />
            <button
              type="submit"
              style={{
                padding: '0.75rem 1rem',
                borderRadius: '0.5rem',
                border: '1px solid #00ff00',
                background: 'black',
                color: '#00ff00',
                fontWeight: 'bold',
                fontSize: '1rem',
                cursor: 'pointer',
                marginTop: '0.5rem',
                textShadow: '0 0 10px #00ff00',
                transition: 'all 0.3s',
              }}
              onMouseOver={(e) => {
                e.currentTarget.style.background = '#00ff00';
                e.currentTarget.style.color = 'black';
              }}
              onMouseOut={(e) => {
                e.currentTarget.style.background = 'black';
                e.currentTarget.style.color = '#00ff00';
              }}
            >
              {t('loginButton')}
            </button>
          </form>
        </div>
      </div>
    </>
  );
}
