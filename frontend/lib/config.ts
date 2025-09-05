export const config = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL!,
};

if (!config.apiUrl) {
  //throw new Error("❌ NEXT_PUBLIC_API_URL is not defined in .env.local");
}
