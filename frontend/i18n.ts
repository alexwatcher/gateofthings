export const locales = ["en", "ru"] as const;
export type Locale = (typeof locales)[number];

export const messages: Record<Locale, () => Promise<any>> = {
  en: () => import("./messages/en.json"),
  ru: () => import("./messages/ru.json"),
};
