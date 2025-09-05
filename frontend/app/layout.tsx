export const dynamic = "force-dynamic";

import "./globals.css";
import MatrixRain from "../components/MatrixRain";
import { LocaleProvider } from "../context/LocaleContext";
import { NotificationProvider } from "../context/NotificationContext";
import { ConfigProvider } from "@/context/ConfigContext";
import { loadConfig } from "../lib/loadConfig";

const availableLocales = ["en", "ru"];

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const config = loadConfig();
  return (
    <html>
      <body className="h-screen w-screen">
        <MatrixRain />
        <LocaleProvider defaultLocale="en" localesList={availableLocales}>
          <NotificationProvider>
            <ConfigProvider value={config}>{children}</ConfigProvider>
          </NotificationProvider>
        </LocaleProvider>
      </body>
    </html>
  );
}
