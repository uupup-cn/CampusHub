import type { Metadata } from 'next';
import './globals.css';
import AuroraBackground from '@/components/AuroraBackground';
import Navbar from '@/components/Navbar';
import { ToastProvider } from '@/components/Toast';

export const metadata: Metadata = {
  title: 'CampusHub — 知识共享积分生态',
  description: '以知识共享为核心的论坛积分生态：发帖、回复、签到赚取 CHB，积分集市流转，OAuth2 应用接入。',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <head>
        <meta charSet="utf-8" />
      </head>
      <body>
        <AuroraBackground />
        <div className="noise-overlay" aria-hidden />
        <ToastProvider>
          <Navbar />
          <main className="pt-16 min-h-screen">{children}</main>
        </ToastProvider>
      </body>
    </html>
  );
}
