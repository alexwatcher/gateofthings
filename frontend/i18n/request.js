import { getRequestConfig } from "next-intl/server";

export default getRequestConfig(() => ({
  locale: "en",
  messages: {}
}));