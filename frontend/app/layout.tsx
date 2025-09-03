import { NextIntlClientProvider } from "next-intl";
import "./globals.css";
import MatrixRain from "../components/MatrixRain";
import { LocaleProvider } from "../context/LocaleContext";

const availableLocales = ["en", "ru"];

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html>
      <body className="h-screen w-screen">
        <MatrixRain />
        <LocaleProvider defaultLocale="en" localesList={availableLocales}>
          {children}
        </LocaleProvider>
      </body>
    </html>
  );
}
