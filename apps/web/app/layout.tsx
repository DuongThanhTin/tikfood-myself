import type { Metadata } from "next";
import "maplibre-gl/dist/maplibre-gl.css";
import "./globals.css";

import { AuthProvider } from "../components/auth/AuthProvider";

export const metadata: Metadata = {
  title: "TikFood Discovery",
  description: "Realtime social food discovery for trending dishes and venues."
};

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
