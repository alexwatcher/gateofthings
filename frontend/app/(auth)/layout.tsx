"use client";

import MatrixRain from "../components/matrixrain";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen items-center justify-center font-mono relative">
      <MatrixRain />
      <div className="relative z-10">{children}</div>
    </div>
  );
}
