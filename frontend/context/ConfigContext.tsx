"use client";

import { createContext, useContext, ReactNode } from "react";
import { AppConfig } from "../types/config";

const ConfigContext = createContext<AppConfig | null>(null);

export function ConfigProvider({
  value,
  children,
}: {
  value: AppConfig;
  children: ReactNode;
}) {
  return (
    <ConfigContext.Provider value={value}>{children}</ConfigContext.Provider>
  );
}

export function useConfig() {
  const ctx = useContext(ConfigContext);
  if (!ctx) throw new Error("useConfig must be used within ConfigProvider");
  return ctx;
}
