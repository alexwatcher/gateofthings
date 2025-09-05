export const dynamic = "force-dynamic";

import { AppConfig } from "../types/config";

export function loadConfig(): AppConfig {
  return {
    apiUrl: process.env.API_URL || "http://localhost:3000",
  };
}
