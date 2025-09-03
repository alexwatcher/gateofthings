"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";
import { NextIntlClientProvider } from "next-intl";
import Cookies from "js-cookie";

type LocaleContextProps = {
  locale: string;
  setLocale: (locale: string) => void;
  availableLocales: readonly string[];
};

const LocaleContext = createContext<LocaleContextProps>({
  locale: "en",
  setLocale: () => {},
  availableLocales: ["en"],
});

export const LocaleProvider = ({
  children,
  defaultLocale = "en",
  localesList = ["en"],
}: {
  children: ReactNode;
  defaultLocale?: string;
  localesList?: readonly string[];
}) => {
  const [locale, setLocaleState] = useState<string>(defaultLocale);
  const [messages, setMessages] = useState<Record<string, string>>({});

  useEffect(() => {
    const saved = Cookies.get("NEXT_LOCALE") || defaultLocale;
    changeLocale(saved);
  }, [defaultLocale]);

  const changeLocale = async (newLocale: string) => {
    console.log("[LocaleProvider] Requested locale change:", newLocale);

    if (!localesList.includes(newLocale)) {
      console.warn(
        "[LocaleProvider] Locale not in available list:",
        newLocale,
        localesList
      );
      return;
    }

    setLocaleState(newLocale);
    Cookies.set("NEXT_LOCALE", newLocale, { expires: 365 });
    console.log(
      "[LocaleProvider] Locale cookie set:",
      Cookies.get("NEXT_LOCALE")
    );

    try {
      console.log("[LocaleProvider] Importing messages for:", newLocale);
      const imported = await import(`../messages/${newLocale}.json`);
      console.log("[LocaleProvider] Messages loaded:", imported);

      setMessages(imported.default);
      console.log("[LocaleProvider] Messages state updated:", imported.default);
    } catch (e) {
      console.error(
        `[LocaleProvider] Cannot load messages for locale ${newLocale}`,
        e
      );
      setMessages({});
    }
  };

  return (
    <LocaleContext.Provider
      value={{ locale, setLocale: changeLocale, availableLocales: localesList }}
    >
      {messages && Object.keys(messages).length > 0 && (
        <NextIntlClientProvider locale={locale} messages={messages}>
          {children}
        </NextIntlClientProvider>
      )}
    </LocaleContext.Provider>
  );
};

export const useLocale = () => useContext(LocaleContext);
