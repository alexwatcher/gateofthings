"use client";

import { useRouter } from "next/navigation";
import { Fragment, useEffect, useState } from "react";
import { Listbox, Transition } from "@headlessui/react";
import Cookies from "js-cookie";
import { locales, Locale } from "../../i18n";

export default function LanguageSwitcher({
  defaultLocale = "en",
}: {
  defaultLocale?: Locale;
}) {
  const router = useRouter();
  const [selected, setSelected] = useState<Locale>(defaultLocale);
  const [displayNames, setDisplayNames] = useState<Record<Locale, string>>(
    {} as Record<Locale, string>
  );
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const dn = new Intl.DisplayNames([defaultLocale], { type: "language" });
    const names: Record<Locale, string> = {} as Record<Locale, string>;
    locales.forEach((code) => {
      names[code] = dn.of(code) || code;
    });
    setDisplayNames(names);

    const saved = Cookies.get("NEXT_LOCALE") as Locale;
    if (saved && locales.includes(saved)) setSelected(saved);
    router.refresh();
  }, [defaultLocale]);

  const handleChange = (locale: Locale) => {
    setSelected(locale);
    Cookies.set("NEXT_LOCALE", locale, { expires: 365 });
    setOpen(false);
  };

  const longestOption = Math.max(
    ...Object.values(displayNames).map((name) => name.length)
  );
  const width = `${longestOption + 2}ch`;

  return (
    <div className="relative z-20" style={{ width }}>
      <div
        className="bg-black text-green-400 font-mono py-2 rounded-md border border-green-500 drop-shadow-[0_0_10px_#00ff00] cursor-pointer flex justify-center"
        onClick={() => setOpen(!open)}
      >
        {displayNames[selected]}
      </div>

      <Transition
        as={Fragment}
        show={open}
        leave="transition ease-in duration-100"
        leaveFrom="opacity-100"
        leaveTo="opacity-0"
      >
        <div className="absolute mt-1 w-full bg-black border border-green-500 rounded-md shadow-lg max-h-60 overflow-auto focus:outline-none z-30 list-none">
          {locales.map((code) => (
            <div
              key={code}
              className={`cursor-pointer select-none relative flex items-center justify-center py-2 font-mono text-green-400 hover:bg-green-900 hover:shadow-[0_0_10px_#00ff00] ${
                code === selected ? "underline" : ""
              }`}
              onClick={() => handleChange(code)}
            >
              {displayNames[code]}
            </div>
          ))}
        </div>
      </Transition>
    </div>
  );
}
