"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import Link from "next/link";

export default function SignInPage() {
  const t = useTranslations();
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSignIn = (e: React.FormEvent) => {
    e.preventDefault();
    router.push("/dashboard");
  };

  return (
    <form
      onSubmit={handleSignIn}
      className="flex w-80 flex-col gap-4 rounded-xl border border-green-500 bg-black/80 p-6 shadow-[0_0_20px_#00ff00]"
    >
      <h1 className="mb-2 text-center text-3xl font-bold text-green-400 tracking-wider drop-shadow-[0_0_10px_#00ff00]">
        {t("appTitle")}
      </h1>

      <input
        type="text"
        placeholder={t("loginPlaceholder")}
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        className="rounded-lg border border-green-500 bg-black p-2 text-green-400 placeholder-green-700 focus:outline-none focus:shadow-[0_0_10px_#00ff00]"
      />

      <input
        type="password"
        placeholder={t("passwordPlaceholder")}
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        className="rounded-lg border border-green-500 bg-black p-2 text-green-400 placeholder-green-700 focus:outline-none focus:shadow-[0_0_10px_#00ff00]"
      />

      <button
        type="submit"
        className="mt-2 rounded-lg border border-green-500 bg-black p-2 font-bold text-green-400 transition hover:bg-green-500 hover:text-black hover:shadow-[0_0_15px_#00ff00]"
      >
        {t("signInButton")}
      </button>

      <div className="mt-3 flex justify-between w-full text-sm">
        <Link
          href="/signup"
          className="text-green-400 hover:text-green-200 transition"
        >
          {t("signUpLink")}
        </Link>
        <Link
          href="/forgot-password"
          className="text-green-400 hover:text-green-200 transition"
        >
          {t("forgotPasswordLink")}
        </Link>
      </div>
    </form>
  );
}
