'use client';

import { useTranslations } from 'next-intl';
import React, { useRef, useLayoutEffect, useState } from 'react';

export default function LoginPage() {
  const t = useTranslations();
  const formRef = useRef<HTMLDivElement>(null);
  const [titleTop, setTitleTop] = useState<number>(0);

  useLayoutEffect(() => {
    if (formRef.current) {
      const formRect = formRef.current.getBoundingClientRect();
      // 75% of the distance from top of screen to top of form
      setTitleTop(formRect.top * 0.50);
    }
  }, []);

  return (
    <div
      style={{
        minHeight: '100vh',
        width: '100vw',
        position: 'relative',
        background: 'linear-gradient(135deg, #1e293b 0%, #64748b 100%)',
        overflow: 'hidden',
      }}
    >
      {/* Title at 75% of the way from top to top of form */}
      <h1
        style={{
          position: 'absolute',
          left: '50%',
          transform: 'translateX(-50%)',
          top: titleTop,
          color: '#fff',
          fontSize: '2.5rem',
          fontWeight: 'bold',
          letterSpacing: '2px',
          textShadow: '0 2px 16px #0006',
          margin: 0,
          zIndex: 3,
          pointerEvents: 'none',
          transition: 'top 0.2s',
        }}
      >
        {t('title')}
      </h1>
      {/* Centered login form */}
      <div
        ref={formRef}
        style={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          zIndex: 2,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <form
          style={{
            background: '#fff',
            padding: '2rem 2.5rem',
            borderRadius: '1rem',
            boxShadow: '0 4px 32px #0002',
            display: 'flex',
            flexDirection: 'column',
            gap: '1.2rem',
            minWidth: '320px',
          }}
        >
          <input
            type="text"
            placeholder={t('loginPlaceholder')}
            style={{
              padding: '0.75rem 1rem',
              borderRadius: '0.5rem',
              border: '1px solid #cbd5e1',
              fontSize: '1rem',
              color: '#222',
              background: '#fff',
            }}
          />
          <input
            type="password"
            placeholder={t('passwordPlaceholder')}
            style={{
              padding: '0.75rem 1rem',
              borderRadius: '0.5rem',
              border: '1px solid #cbd5e1',
              fontSize: '1rem',
              color: '#222',
              background: '#fff',
            }}
          />
          <button
            type="submit"
            style={{
              padding: '0.75rem 1rem',
              borderRadius: '0.5rem',
              border: 'none',
              background: '#334155',
              color: '#fff',
              fontWeight: 'bold',
              fontSize: '1rem',
              cursor: 'pointer',
              marginTop: '0.5rem',
            }}
          >
            {t('loginButton')}
          </button>
        </form>
      </div>
    </div>
  );
}