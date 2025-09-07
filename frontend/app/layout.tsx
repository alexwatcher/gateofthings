// export const dynamic = "force-dynamic";
import "./globals.css";
import MatrixRain from "../components/MatrixRain";
import { LocaleProvider } from "../context/LocaleContext";
import { NotificationProvider } from "../context/NotificationContext";
import { ConfigProvider } from "@/context/ConfigContext";

const availableLocales = ["en", "ru"];

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const varName = "API_URL";
  console.info("API_URL: ", process.env[varName]);

  // const config = loadConfig();
  return (
    <html>
      <body className="h-screen w-screen">
        <MatrixRain />
        <LocaleProvider defaultLocale="en" localesList={availableLocales}>
          <NotificationProvider>{children}</NotificationProvider>
        </LocaleProvider>
      </body>
    </html>
  );
}
