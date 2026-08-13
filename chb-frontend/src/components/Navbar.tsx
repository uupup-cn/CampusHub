'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useState, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Store, LayoutDashboard, ShieldCheck, Sparkles,
  MessageSquare, LogIn, LogOut, Wallet, User, ChevronDown,
  Package, ShoppingBag, PlusCircle, Settings, Clock,
} from 'lucide-react';
import { apiRequest, formatCHB } from '@/lib/api';

const FORUM_URL = process.env.NEXT_PUBLIC_DISCOURSE_URL || 'http://198.44.177.228:9800';

interface UserInfo {
  user_id: number;
  username: string;
  trust_level: number;
  avatar_template?: string;
}

function getAvatarUrl(avatarTemplate?: string): string {
  if (!avatarTemplate) return '';
  const url = avatarTemplate.replace('{size}', '40');
  return url.startsWith('http') ? url : FORUM_URL + url;
}

export default function Navbar() {
  const pathname = usePathname();
  const isAdmin = pathname?.startsWith('/admin');
  const isDashboard = pathname?.startsWith('/dashboard');
  const isMarketplace = pathname?.startsWith('/marketplace');
  const [user, setUser] = useState<UserInfo | null>(null);
  const [balance, setBalance] = useState<{ balance: number; pending_balance: number } | null>(null);
  const [merchantStatus, setMerchantStatus] = useState<{ is_merchant: boolean; status: string } | null>(null);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const hoverRef = useRef<HTMLDivElement>(null);
  const closeTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const fetchAuth = () => {
    fetch('/api/auth/me').then(r => r.json()).then(data => {
      if (data.code === 0 && data.data?.logged_in) {
        const u = data.data;
        setUser({ user_id: u.user_id, username: u.username, trust_level: u.trust_level, avatar_template: u.avatar_template });
        localStorage.setItem('chb_user_id', String(u.user_id));
        return apiRequest<{ balance: number; pending_balance: number }>('/api/chb/balance');
      }
      return null;
    }).then(res => {
      if (res && res.code === 0) setBalance(res.data);
    }).catch(() => {});
  };

  useEffect(() => {
    fetchAuth();
    const handleAuthChange = () => {
      setTimeout(fetchAuth, 300);
    };
    window.addEventListener('auth-state-changed', handleAuthChange);
    return () => window.removeEventListener('auth-state-changed', handleAuthChange);
  }, []);

  // Fetch merchant status when on marketplace page
  useEffect(() => {
    if (isMarketplace && user) {
      apiRequest<{ is_merchant: boolean; status: string }>('/api/marketplace/my-status')
        .then(res => { if (res.code === 0) setMerchantStatus(res.data); })
        .catch(() => {});
    }
  }, [isMarketplace, user]);

  const handleMouseEnter = () => {
    if (closeTimer.current) clearTimeout(closeTimer.current);
    setDropdownOpen(true);
  };

  const handleMouseLeave = () => {
    closeTimer.current = setTimeout(() => setDropdownOpen(false), 200);
  };

  const handleLogout = () => {
    window.location.href = FORUM_URL + '/logout?return_url=' + window.location.origin;
  };

  const links = isAdmin
    ? [
        { href: '/admin', label: '仪表盘', icon: LayoutDashboard },
        { href: '/admin/marketplace', label: '集市审核', icon: Store },
        { href: '/admin/users', label: '用户管理', icon: ShieldCheck },
      ]
    : isDashboard
    ? [
        { href: '/dashboard', label: '概览', icon: LayoutDashboard },
        { href: '/dashboard/marketplace', label: '集市', icon: Store },
        { href: '/dashboard/points', label: '积分', icon: Wallet },
        { href: '/dashboard/apps', label: '授权应用', icon: Package },
      ]
    : [
        { href: '/marketplace', label: '集市', icon: Store },
      ];

  // Build dropdown menu items based on context
  const dropdownItems: Array<{ label: string; href?: string; onClick?: () => void; icon: typeof User }> = [];
  dropdownItems.push({ label: '个人中心', href: '/dashboard', icon: LayoutDashboard });
  if (isMarketplace) {
    if (merchantStatus?.is_merchant) {
      dropdownItems.push({ label: '我的商品', href: '/dashboard/marketplace', icon: Package });
      dropdownItems.push({ label: '我的订单', href: '/dashboard/marketplace/orders', icon: ShoppingBag });
    } else {
      dropdownItems.push({ label: '申请入驻', href: '/dashboard/marketplace/apply', icon: PlusCircle });
    }
  }

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-white/5 backdrop-blur-xl bg-[#050508]/70">
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2 group">
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

          {user ? (
            <div
              ref={hoverRef}
              className="relative"
              onMouseEnter={handleMouseEnter}
              onMouseLeave={handleMouseLeave}
            >
              <button className="flex items-center gap-2 pl-1 pr-2 py-1 rounded-full border border-white/10 hover:border-white/25 transition-all">
                {user.avatar_template ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={getAvatarUrl(user.avatar_template)}
                    alt={user.username}
                    className="w-7 h-7 rounded-full object-cover"
                  />
                ) : (
                  <span className="w-7 h-7 rounded-full bg-gradient-to-br from-violet-500 to-cyan-400 flex items-center justify-center">
                    <User size={14} className="text-white" />
                  </span>
                )}
                <span className="text-xs font-medium text-gray-300 max-w-[80px] truncate">{user.username}</span>
                <span className="badge badge-purple !text-[10px] !px-1.5 !py-0.5">TL{user.trust_level}</span>
                <ChevronDown size={14} className={'text-gray-500 transition-transform ' + (dropdownOpen ? 'rotate-180' : '')} />
              </button>

              <AnimatePresence>
                {dropdownOpen && (
                  <motion.div
                    initial={{ opacity: 0, y: -8 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -8 }}
                    transition={{ duration: 0.2, ease: [0.22, 1, 0.36, 1] }}
                    className="absolute right-0 top-full mt-2 w-56 glass-card !rounded-xl overflow-hidden py-1"
                  >
                    {balance !== null && (
                      <div className="px-4 py-2.5 border-b border-white/5 space-y-1">
                        <div className="flex items-center gap-2">
                          <Wallet size={14} className="text-amber-300" />
                          <span className="text-xs text-gray-400">可用积分</span>
                          <span className="text-sm font-bold text-amber-300 ml-auto">{formatCHB(balance.balance)}</span>
                        </div>
                        {balance.pending_balance > 0 && (
                          <div className="flex items-center gap-2">
                            <Clock size={14} className="text-violet-300" />
                            <span className="text-xs text-gray-400">未来积分</span>
                            <span className="text-sm font-bold text-violet-300 ml-auto">{formatCHB(balance.pending_balance)}</span>
                          </div>
                        )}
                      </div>
                    )}
                    {dropdownItems.map((item, i) => {
                      const Icon = item.icon;
                      return item.href ? (
                        <Link
                          key={i}
                          href={item.href}
                          className="flex items-center gap-2 px-4 py-2.5 text-sm text-gray-300 hover:bg-white/5 hover:text-white transition-colors"
                          onClick={() => setDropdownOpen(false)}
                        >
                          <Icon size={15} className="text-gray-500" />
                          {item.label}
                        </Link>
                      ) : (
                        <button
                          key={i}
                          onClick={item.onClick}
                          className="flex items-center gap-2 px-4 py-2.5 text-sm text-gray-300 hover:bg-white/5 hover:text-white transition-colors w-full text-left"
                        >
                          <Icon size={15} className="text-gray-500" />
                          {item.label}
                        </button>
                      );
                    })}
                    <div className="border-t border-white/5 mt-1 pt-1">
                      <button
                        onClick={handleLogout}
                        className="flex items-center gap-2 px-4 py-2.5 text-sm text-gray-400 hover:bg-red-500/10 hover:text-red-400 transition-colors w-full text-left"
                      >
                        <LogOut size={15} />
                        退出登录
                      </button>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
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
