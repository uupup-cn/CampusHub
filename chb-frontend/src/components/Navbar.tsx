'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import { Store, LayoutDashboard, ShieldCheck, Sparkles, MessageSquare, LogIn, LogOut, Wallet, User } from 'lucide-react';
import { apiRequest, formatCHB } from '@/lib/api';

const FORUM_URL = process.env.NEXT_PUBLIC_FORUM_URL || 'http://112.213.106.104:9800';

interface UserInfo {
  user_id: number;
  username: string;
  trust_level: number;
}

export default function Navbar() {
  const pathname = usePathname();
  const isAdmin = pathname?.startsWith('/admin');
  const [user, setUser] = useState<UserInfo | null>(null);
  const [balance, setBalance] = useState<number | null>(null);

  useEffect(() => {
    fetch('/api/auth/me').then(r => r.json()).then(data => {
      if (data.code === 0 && data.data?.logged_in) {
        const u = data.data;
        setUser({ user_id: u.user_id, username: u.username, trust_level: u.trust_level });
        localStorage.setItem('chb_user_id', String(u.user_id));
        // 设置 user_id 后立即查余额
        return apiRequest<{ balance: number }>('/api/chb/balance');
      }
      return null;
    }).then(res => {
      if (res && res.code === 0) setBalance(res.data.balance);
    }).catch(() => {});
  }, []);

  const handleLogout = () => {
    window.location.href = FORUM_URL + '/logout?return_url=' + window.location.origin;
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
                className={'nav-link flex items-center gap-1.5 ' + (active ? '!text-white' : '')}
              >
                <Icon size={15} />
                {l.label}
              </Link>
            );
          })}

          <a
            href={FORUM_URL}
            className="nav-link flex items-center gap-1.5"
            target="_blank"
            rel="noopener noreferrer"
          >
            <MessageSquare size={15} />
            论坛
          </a>

          {isAdmin && (
            <Link href="/marketplace" className="nav-link flex items-center gap-1.5">
              <Store size={15} />
              返回集市
            </Link>
          )}

          {user ? (
            <div className="flex items-center gap-3">
              <span className="flex items-center gap-1.5 text-xs text-gray-400">
                <User size={14} className="text-cyan-400" />
                {user.username}
                <span className="badge badge-purple !text-[10px] !px-1.5 !py-0.5">TL{user.trust_level}</span>
              </span>
              {balance !== null && (
                <span className="flex items-center gap-1 text-xs text-amber-300">
                  <Wallet size={13} />
                  {formatCHB(balance)}
                </span>
              )}
              <button
                onClick={handleLogout}
                className="flex items-center gap-1.5 px-2 py-1.5 rounded-full border border-white/10 text-xs text-gray-500 hover:text-red-400 hover:border-red-500/30 transition-all"
                title="退出登录"
              >
                <LogOut size={13} />
              </button>
            </div>
          ) : (
            <button
              onClick={() => { window.location.href = FORUM_URL + '/login?return_url=' + window.location.origin + '/auth/callback'; }}
              className="flex items-center gap-2 px-3 py-1.5 rounded-full border border-white/10 text-xs text-gray-400 hover:text-white hover:border-white/25 transition-all"
              title="登录论坛账号"
            >
              <LogIn size={14} className="text-violet-400" />
              登录
            </button>
          )}
        </div>
      </div>
    </nav>
  );
}
