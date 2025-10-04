import Cookies from "js-cookie";

type FetchOptions = {
  body?: any;
  headers?: Record<string, string>;
};

export async function postWithCsrf<T = any>(
  url: string,
  options: FetchOptions = {}
): Promise<T> {
  const csrfToken = Cookies.get("csft");

  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      csrf: csrfToken || "",
      ...options.headers,
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const data = await res.json().catch(() => null);

  if (!res.ok) {
    throw new Error(data?.message || `HTTP error ${res.status}`);
  }

  return data;
}
