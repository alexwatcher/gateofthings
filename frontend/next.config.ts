import type { NextConfig } from "next";

const createNextIntlPlugin = require("next-intl/plugin");
const withNextIntl = createNextIntlPlugin({});

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
};

export default withNextIntl(nextConfig);
