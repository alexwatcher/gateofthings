"use client";

import { Fragment, useEffect, useRef, useState } from "react";
import { Transition } from "@headlessui/react";
import { useLocale } from "../context/LocaleContext";

export default function LanguageSwitcher() {
  const { locale, setLocale, availableLocales } = useLocale();
  const [displayNames, setDisplayNames] = useState<Record<string, string>>({});
  const [open, setOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Prepare display names dynamically using Intl.DisplayNames
  useEffect(() => {
    const dn = new Intl.DisplayNames(["en"], { type: "language" });
    const names: Record<string, string> = {};
    availableLocales.forEach((code) => {
      names[code] = dn.of(code) || code;
    });
    setDisplayNames(names);
  }, [locale, availableLocales]);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleChange = (newLocale: string) => {
    setLocale(newLocale);
    setOpen(false);
  };

  const longestOption = Math.max(
    ...Object.values(displayNames).map((name) => name.length)
  );
  const width = `${longestOption + 2}ch`;

  return (
    <div className="relative z-20" style={{ width }} ref={dropdownRef}>
      {/* Button showing current language */}
      <div
        className="bg-black text-green-400 font-mono py-2 rounded-md border border-green-500 drop-shadow-[0_0_10px_#00ff00] cursor-pointer flex justify-center"
        onClick={() => setOpen(!open)}
      >
        {displayNames[locale]}
      </div>

      {/* Dropdown list */}
      <Transition
        as={Fragment}
        show={open}
        leave="transition ease-in duration-100"
        leaveFrom="opacity-100"
        leaveTo="opacity-0"
      >
        <div className="absolute mt-1 w-full bg-black border border-green-500 rounded-md shadow-lg max-h-60 overflow-auto focus:outline-none z-30 list-none">
          {availableLocales.map((code) => (
            <div
              key={code}
              className={`cursor-pointer select-none relative flex items-center justify-center py-2 font-mono text-green-400 hover:bg-green-900 hover:shadow-[0_0_10px_#00ff00] ${
                code === locale ? "underline" : ""
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
