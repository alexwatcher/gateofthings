"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useNotification } from "@/context/NotificationContext";
// import { useConfig } from "@/context/ConfigContext";

export default function SignupPagePage() {
  const router = useRouter();
  const t = useTranslations();
  const { notify } = useNotification();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  // const config = useConfig();

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await fetch(`{config.apiUrl}/auth/signup`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        throw new Error(`HTTP error: ${res.status}`);
      }

      const data = await res.json();
      localStorage.setItem("token", data.token);

      notify(t("signupSuccess"), "success");
      router.back();
    } catch (err) {
      console.error("Signup failed:", err);
      notify(t("signupError"), "error");
    }
  };

  return (
    <form
      onSubmit={handleSignUp}
      className="relative flex w-80 flex-col gap-4 rounded-xl border border-green-500 bg-black/80 p-6 shadow-[0_0_20px_#00ff00]"
    >
      <button
        type="button"
        onClick={() => router.back()}
        className="absolute top-1.5 right-3 text-green-400 hover:text-red-500 font-bold text-xl"
      >
        ×
      </button>

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

      <input
        type="password"
        placeholder={t("confirmPasswordPlaceholder")}
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
        className="rounded-lg border border-green-500 bg-black p-2 text-green-400 placeholder-green-700 focus:outline-none focus:shadow-[0_0_10px_#00ff00]"
      />

      <button
        type="submit"
        className="mt-2 rounded-lg border border-green-500 bg-black p-2 font-bold text-green-400 transition hover:bg-green-500 hover:text-black hover:shadow-[0_0_15px_#00ff00]"
      >
        {t("signUpButton")}
      </button>
    </form>
  );
}
