"use client";

import { createContext, useContext, useState, ReactNode } from "react";

type Notification = {
  id: number;
  message: string;
  type: "error" | "warning" | "info";
};

type NotificationContextType = {
  notify: (msg: string, type?: "error" | "warning" | "info") => void;
};

const NotificationContext = createContext<NotificationContextType | undefined>(
  undefined
);

export function NotificationProvider({ children }: { children: ReactNode }) {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const notify = (
    msg: string,
    type: "error" | "warning" | "info" = "error"
  ) => {
    const id = Date.now();
    const n: Notification = { id, message: msg, type };

    setNotifications((prev) => [n, ...prev]);

    setTimeout(() => {
      setNotifications((prev) => prev.filter((x) => x.id !== id));
    }, 2500);
  };

  return (
    <NotificationContext.Provider value={{ notify }}>
      {children}

      {/* Notification stack */}
      <div className="fixed top-4 left-1/2 -translate-x-1/2 z-50 flex flex-col gap-3 items-center">
        {notifications.map((notification) => (
          <div
            key={notification.id}
            className={`rounded-md px-6 py-3 font-mono text-sm tracking-widest shadow-[0_0_20px] animate-matrix-fade
              ${
                notification.type === "error"
                  ? "bg-black text-red-500 shadow-red-500"
                  : notification.type === "warning"
                  ? "bg-black text-yellow-400 shadow-yellow-400"
                  : "bg-black text-green-400 shadow-green-400"
              }`}
            style={{
              textShadow:
                notification.type === "error"
                  ? "0 0 10px #ff0000, 0 0 20px #ff0000, 0 0 30px #ff0000"
                  : notification.type === "warning"
                  ? "0 0 10px #ffff00, 0 0 20px #ffff00, 0 0 30px #ffff00"
                  : "0 0 10px #00ff00, 0 0 20px #00ff00, 0 0 30px #00ff00",
            }}
          >
            {notification.message.split(";").map((line, i) => (
              <div key={i}>{line.trim()}</div>
            ))}
          </div>
        ))}
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
    </NotificationContext.Provider>
  );
}

export const useNotification = () => {
  const ctx = useContext(NotificationContext);
  if (!ctx)
    throw new Error("useNotification must be used inside NotificationProvider");
  return ctx;
};
