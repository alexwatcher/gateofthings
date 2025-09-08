"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useNotification } from "@/context/NotificationContext";
import Link from "next/link";

export default function SignInPage() {
  const t = useTranslations();
  const router = useRouter();
  const { notify } = useNotification();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSignIn = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/auth/signin`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, password }),
        }
      );

      if (!res.ok) {
        throw new Error(`HTTP error: ${res.status}`);
      }

      const data = await res.json();
      localStorage.setItem("token", data.token);

      notify(t("signinSuccess"), "success");
      router.push("/dashboard");
    } catch (err) {
      console.error("Signin failed:", err);
      notify(t("signinError"), "error");
    }
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
        placeholder={t("emailPlaceholder")}
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
          href="/resetpassword"
          className="text-green-400 hover:text-green-200 transition"
        >
          {t("restorePasswordLink")}
        </Link>
      </div>
    </form>
  );
}
