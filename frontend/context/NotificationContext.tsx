// context/notificationContext.tsx
"use client";

import { createContext, useContext, useState, ReactNode } from "react";

type Notification = {
  message: string;
  type?: "error" | "success";
};

type NotificationContextType = {
  notify: (msg: string, type?: "error" | "success") => void;
};

const NotificationContext = createContext<NotificationContextType | undefined>(
  undefined
);

export function NotificationProvider({ children }: { children: ReactNode }) {
  const [notification, setNotification] = useState<Notification | null>(null);

  const notify = (msg: string, type: "error" | "success" = "error") => {
    setNotification({ message: msg, type });
    // total duration 2.5s = 0.25s fade-in + 2s visible + 0.25s fade-out
    setTimeout(() => setNotification(null), 2500);
  };

  return (
    <NotificationContext.Provider value={{ notify }}>
      {children}

      {notification && (
        <>
          {/* Overlay to block input */}
          <div className="fixed inset-0 z-40 cursor-not-allowed" />

          {/* Floating notification */}
          <div
            className={`fixed top-4 left-1/2 -translate-x-1/2 z-50 rounded-md px-6 py-3 font-mono text-lg tracking-widest shadow-[0_0_20px]
              ${
                notification.type === "error"
                  ? "bg-black text-red-500 shadow-red-500"
                  : "bg-black text-green-400 shadow-green-400"
              } animate-matrix-fade`}
            style={{
              textShadow:
                notification.type === "error"
                  ? "0 0 10px #ff0000, 0 0 20px #ff0000, 0 0 30px #ff0000"
                  : "0 0 10px #00ff00, 0 0 20px #00ff00, 0 0 30px #00ff00",
            }}
          >
            {notification.message}
          </div>

          {/* Inline animation CSS */}
          <style jsx>{`
            @keyframes matrix-fade {
              0% {
                opacity: 0;
                transform: translateY(-20px) scale(1);
                filter: blur(2px);
              }
              10% {
                opacity: 1;
                transform: translateY(0) scale(1.05);
                filter: blur(0);
              }
              90% {
                opacity: 1;
                transform: translateY(0) scale(1.05);
                filter: blur(0);
              }
              100% {
                opacity: 0;
                transform: translateY(-10px) scale(1);
                filter: blur(1px);
              }
            }

            .animate-matrix-fade {
              animation: matrix-fade 2.5s ease-in-out forwards;
            }
          `}</style>
        </>
      )}
    </NotificationContext.Provider>
  );
}

export const useNotification = () => {
  const ctx = useContext(NotificationContext);
  if (!ctx)
    throw new Error("useNotification must be used inside NotificationProvider");
  return ctx;
};
