import type { Metadata } from "next";
import "./globals.css";

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
      <body>{children}</body>
    </html>
  );
}
