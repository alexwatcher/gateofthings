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
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      if (confirm !== password) {
        notify(t("passwordsDontMatch"), "error");
        return;
      }

      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/v1/auth/signup`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email: email, password: password }),
        }
      );

      if (!res.ok) {
        let errorMsg = t("signupError");
        try {
          const errData = await res.json();
          if (errData?.message) {
            errorMsg = errData.message;
          }
        } catch {}
        throw new Error(errorMsg);
      }

      notify(t("signupSuccess"), "info");
      router.back();
    } catch (err: unknown) {
      const message =
        err instanceof Error
          ? err.message
          : typeof err === "string"
          ? err
          : t("signinError");
      notify(message, "error");
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
        value={email}
        onChange={(e) => setEmail(e.target.value)}
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
