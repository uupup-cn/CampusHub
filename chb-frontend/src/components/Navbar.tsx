'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Store, LayoutDashboard, ShieldCheck, Sparkles } from 'lucide-react';
import { useState } from 'react';
import { getUserID, setUserID } from '@/lib/api';

export default function Navbar() {
  const pathname = usePathname();
  const [userId, setUserIdState] = useState(getUserID());
  const isAdmin = pathname?.startsWith('/admin');

  const switchUser = () => {
    const next = userId === 1 ? 2 : userId === 2 ? 1 : 2;
    setUserID(next);
    setUserIdState(next);
    window.location.reload();
  };

  const links = isAdmin
    ? [
        { href: '/admin', label: '仪表盘', icon: LayoutDashboard },
        { href: '/admin/marketplace', label: '集市审核', icon: Store },
        { href: '/admin/users', label: '用户管理', icon: ShieldCheck },
      ]
    : [
        { href: '/marketplace', label: '集市', icon: Store },
      ];

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-white/5 backdrop-blur-xl bg-[#050508]/70">
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
        <Link href="/marketplace" className="flex items-center gap-2 group">
          <span className="w-8 h-8 rounded-lg bg-gradient-to-br from-violet-500 to-cyan-400 flex items-center justify-center shadow-lg shadow-violet-500/20 group-hover:scale-110 transition-transform">
            <Sparkles size={16} className="text-white" />
          </span>
          <span className="font-bold text-lg tracking-tight">
            Campus<span className="text-violet-400">Hub</span>
          </span>
        </Link>

        <div className="flex items-center gap-6">
          {links.map((l) => {
            const Icon = l.icon;
            const active = pathname === l.href;
            return (
              <Link
                key={l.href}
                href={l.href}
                className={`nav-link flex items-center gap-1.5 ${active ? '!text-white' : ''}`}
              >
                <Icon size={15} />
                {l.label}
              </Link>
            );
          })}

          {isAdmin && (
            <Link href="/marketplace" className="nav-link flex items-center gap-1.5">
              <Store size={15} />
              返回集市
            </Link>
          )}

          <button
            onClick={switchUser}
            className="flex items-center gap-2 px-3 py-1.5 rounded-full border border-white/10 text-xs text-gray-400 hover:text-white hover:border-white/25 transition-all"
            title="切换模拟用户（开发环境）"
          >
            <span className="w-5 h-5 rounded-full bg-gradient-to-br from-violet-500 to-cyan-400 flex items-center justify-center text-[10px] font-bold text-white">
              U{userId}
            </span>
            {userId === 1 ? '用户 U1' : '用户 U2'}
          </button>
        </div>
      </div>
    </nav>
  );
}