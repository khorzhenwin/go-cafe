import "./globals.css";

export const metadata = {
  title: "go-cafe frontend",
  description: "Frontend app for go-cafe"
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
