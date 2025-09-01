import { NextIntlClientProvider } from 'next-intl';
import '../globals.css';

export default function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { locale: string };
}) {
  // Wrap children with NextIntlClientProvider for localization
  return (
    <html lang={params.locale}>
      <body>
        <NextIntlClientProvider locale={params.locale}>
          {children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}