"use client";

import LanguageSwitcher from "../components/languageswitcher";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen items-center justify-center font-mono relative">
      <div className="absolute top-4 right-4 z-20">
        <LanguageSwitcher />
      </div>
      <div className="relative z-10">{children}</div>
    </div>
  );
}
