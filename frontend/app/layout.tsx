import { NextIntlClientProvider } from "next-intl";
import "./globals.css";
import { getLocale, getMessages } from "next-intl/server";
import MatrixRain from "./components/matrixrain";

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const messages = await getMessages();
  const locale = await getLocale();
  return (
    <html lang={locale}>
      <body className="h-screen w-screen">
        <MatrixRain />
        <NextIntlClientProvider messages={messages}>
          {children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
