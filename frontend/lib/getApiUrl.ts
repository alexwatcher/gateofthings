export const dynamic = "force-dynamic";

import { AppConfig } from "../types/config";

export function getApiUrl(): AppConfig {
  return {
    apiUrl: process.env.API_URL || "http://localhost:3000",
  };
}
