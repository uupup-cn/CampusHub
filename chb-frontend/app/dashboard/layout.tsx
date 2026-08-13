'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  LayoutDashboard, Store, Wallet, Package, LogOut, User, Sparkles,
} from 'lucide-react';
import { apiRequest, formatCHB } from '@/lib/api';

const FORUM_URL = process.env.NEXT_PUBLIC_DISCOURSE_URL || 'http://198.44.177.228:9800';

interface UserInfo {
  username: string;
  trust_level: number;
  avatar_template?: string;
}

function getAvatarUrl(avatarTemplate?: string): string {
  if (!avatarTemplate) return '';
  const url = avatarTemplate.replace('{size}', '48');
  return url.startsWith('http') ? url : FORUM_URL + url;
}

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [user, setUser] = useState<UserInfo | null>(null);
  const [balance, setBalance] = useState<{ balance: number; pending_balance: number } | null>(null);

  useEffect(() => {
    fetch('/api/auth/me').then(r => r.json()).then(data => {
      if (data.code === 0 && data.data?.logged_in) {
        setUser({ username: data.data.username, trust_level: data.data.trust_level, avatar_template: data.data.avatar_template });
        localStorage.setItem('chb_user_id', String(data.data.user_id));
        return apiRequest<{ balance: number; pending_balance: number }>('/api/chb/balance');
      }
      return null;
    }).then(res => {
      if (res && res.code === 0) setBalance(res.data);
    }).catch(() => {});
  }, []);

  const handleLogout = () => {
    window.location.href = FORUM_URL + '/logout?return_url=' + window.location.origin;
  };

  const navItems = [
    { href: '/dashboard', label: '概览', icon: LayoutDashboard },
    { href: '/dashboard/marketplace', label: '集市管理', icon: Store },
    { href: '/dashboard/points', label: '积分明细', icon: Wallet },
    { href: '/dashboard/apps', label: '授权应用', icon: Package },
  ];

  return (
    <div className="min-h-screen flex">
      {/* Sidebar */}
      <aside className="w-60 border-r border-white/5 bg-[#080810] flex flex-col flex-shrink-0">
        {/* User header */}
        <div className="p-4 border-b border-white/5">
          <div className="flex items-center gap-3">
            {user?.avatar_template ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img src={getAvatarUrl(user.avatar_template)} alt={user.username} className="w-10 h-10 rounded-full object-cover" />
            ) : (
              <span className="w-10 h-10 rounded-full bg-gradient-to-br from-violet-500 to-cyan-400 flex items-center justify-center">
                <User size={18} className="text-white" />
              </span>
            )}
            <div className="flex-1 min-w-0">
              <div className="text-sm font-medium truncate">{user?.username || '加载中'}</div>
              <div className="text-xs text-gray-500">Trust Level {user?.trust_level ?? '-'}</div>
            </div>
          </div>
        </div>

        {/* Balance */}
        <div className="px-4 py-3 border-b border-white/5 space-y-2">
          <div>
            <div className="text-xs text-gray-500 mb-1">可用积分</div>
            <div className="text-lg font-bold text-amber-300">
              {balance !== null ? formatCHB(balance.balance) : '加载中'}
            </div>
          </div>
          {balance !== null && balance.pending_balance > 0 && (
            <div>
              <div className="text-xs text-gray-500 mb-1">未来积分</div>
              <div className="text-sm font-bold text-violet-300">
                {formatCHB(balance.pending_balance)}
              </div>
            </div>
          )}
        </div>

        {/* Nav */}
        <nav className="flex-1 py-4 px-2">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors mb-1 ' +
                  (active ? 'bg-violet-500/15 text-violet-300 font-medium' : 'text-gray-400 hover:bg-white/5 hover:text-white')}
              >
                <Icon size={16} />
                {item.label}
              </Link>
            );
          })}
        </nav>

        {/* Logout */}
        <div className="p-2 border-t border-white/5">
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-gray-400 hover:bg-red-500/10 hover:text-red-400 transition-colors w-full"
          >
            <LogOut size={16} />
            退出登录
          </button>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-y-auto">
        <div className="max-w-5xl mx-auto px-8 py-10">
          {children}
        </div>
      </main>
    </div>
  );
}
