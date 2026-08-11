'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import {
  Users, Activity, ShoppingBag, TrendingUp, Wallet,
  Store, Settings, FileText, ArrowUpRight, CircleDollarSign,
  PackageCheck, ShieldAlert, Gauge,
} from 'lucide-react';
import { apiRequest, formatCHB } from '@/lib/api';
import Reveal from '@/components/Reveal';

interface Stats {
  total_users: number;
  active_users_today: number;
  total_transactions: number;
  today_transactions: number;
  total_marketplace_orders: number;
  public_pool_water_level: number;
  pending_applications: number;
  pending_items: number;
}

export default function AdminDashboard() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [pools, setPools] = useState<{ public_pool: { balance: number; total_supply: number }, official_pool: { balance: number } } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      apiRequest<Stats>('/api/admin/stats'),
      apiRequest<typeof pools>('/api/chb/pools'),
    ]).then(([s, p]) => {
      if (s.code === 0) setStats(s.data);
      if (p.code === 0) setPools(p.data);
    }).finally(() => setLoading(false));
  }, []);

  const statCards = [
    { label: '总用户', value: stats?.total_users ?? 0, icon: Users, color: 'text-violet-400', bg: 'bg-violet-500/10' },
    { label: '今日活跃', value: stats?.active_users_today ?? 0, icon: Activity, color: 'text-cyan-400', bg: 'bg-cyan-500/10' },
    { label: '总交易', value: stats?.total_transactions ?? 0, icon: TrendingUp, color: 'text-emerald-400', bg: 'bg-emerald-500/10' },
    { label: '今日交易', value: stats?.today_transactions ?? 0, icon: CircleDollarSign, color: 'text-amber-400', bg: 'bg-amber-500/10' },
    { label: '集市订单', value: stats?.total_marketplace_orders ?? 0, icon: ShoppingBag, color: 'text-pink-400', bg: 'bg-pink-500/10' },
  ];

  const waterLevel = stats?.public_pool_water_level ?? (pools ? pools.public_pool.balance / pools.public_pool.total_supply : 0);

  const navCards = [
    { href: '/admin/marketplace', title: '集市审核', desc: `待审核商品 ${stats?.pending_items ?? 0} · 入驻 ${stats?.pending_applications ?? 0}`, icon: PackageCheck },
    { href: '/admin/users', title: '用户管理', desc: '冻结 / 解冻 / 积分追回', icon: ShieldAlert },
    { href: '/admin/apps', title: '应用管理', desc: 'OAuth2 应用 CRUD', icon: Store },
    { href: '/admin/settings', title: '系统配置', desc: '奖励规则 / 等级上限 / 经济参数', icon: Settings },
    { href: '/admin/audit', title: '审计日志', desc: '操作记录与追踪', icon: FileText },
  ];

  return (
    <div className="max-w-7xl mx-auto px-6 py-16">
      <div className="mb-12">
        <span className="section-label flex items-center gap-2 mb-3">
          <Gauge size={13} />
          Admin Console
        </span>
        <h1 className="display-title">管理后台</h1>
      </div>

      {loading ? (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10">
          {[1,2,3,4].map(i => <div key={i} className="glass-card h-28 shimmer" />)}
        </div>
      ) : (
        <>
          {/* Stats */}
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-10">
            {statCards.map((s) => {
              const Icon = s.icon;
              return (
                <motion.div
                  key={s.label}
                  className="glass-card stat-card"
                  initial={{ opacity: 0, y: 24 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.5, delay: s.label.length * 0.02, ease: [0.22, 1, 0.36, 1] }}
                >
                  <div className={`w-9 h-9 rounded-lg ${s.bg} flex items-center justify-center mb-3`}>
                    <Icon size={17} className={s.color} />
                  </div>
                  <div className="text-xs text-gray-500 mb-1">{s.label}</div>
                  <div className="stat-value !text-xl">{s.value.toLocaleString()}</div>
                </motion.div>
              );
            })}
          </div>

          {/* Pool status */}
          <Reveal className="mb-10">
            <div className="glass-card p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <Wallet size={16} className="text-amber-400" />
                <span className="font-semibold">积分池状态</span>
              </div>
              <span className={`badge ${waterLevel > 0.2 ? 'badge-green' : 'badge-red'}`}>
                水位 {(waterLevel * 100).toFixed(1)}%
              </span>
            </div>
            <div className="grid grid-cols-2 gap-6">
              <div>
                <div className="text-xs text-gray-500 mb-2">公共池</div>
                <div className="font-bold text-lg mb-3">{formatCHB(pools?.public_pool?.balance ?? 0)}</div>
                <div className="h-2 rounded-full bg-white/5 overflow-hidden">
                  <div
                    className="h-full rounded-full bg-gradient-to-r from-violet-500 to-cyan-400 transition-all duration-1000"
                    style={{ width: `${Math.min(100, waterLevel * 100)}%` }}
                  />
                </div>
                <div className="text-xs text-gray-600 mt-2">
                  总量 {formatCHB(pools?.public_pool?.total_supply ?? 0)}
                </div>
              </div>
              <div>
                <div className="text-xs text-gray-500 mb-2">官方池（手续费收入）</div>
                <div className="font-bold text-lg text-amber-300 mb-3">{formatCHB(pools?.official_pool?.balance ?? 0)}</div>
                <div className="h-2 rounded-full bg-white/5 overflow-hidden">
                  <div className="h-full rounded-full bg-gradient-to-r from-amber-500 to-orange-400" style={{ width: '100%' }} />
                </div>
                <div className="text-xs text-gray-600 mt-2">来自集市与应用手续费</div>
              </div>
            </div>
            </div>
          </Reveal>
        </>
      )}

      {/* Navigation */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {navCards.map((card) => {
          const Icon = card.icon;
          return (
            <motion.div
              key={card.href}
              initial={{ opacity: 0, y: 24 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: card.href.length * 0.015, ease: [0.22, 1, 0.36, 1] }}
            >
              <Link href={card.href} className="glass-card p-6 group block">
                <div className="flex items-start justify-between mb-4">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-violet-500/20 to-cyan-500/10 flex items-center justify-center">
                    <Icon size={18} className="text-violet-300" />
                  </div>
                  <ArrowUpRight size={16} className="text-gray-600 group-hover:text-white group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-all" />
                </div>
                <h3 className="font-semibold mb-1 group-hover:text-violet-300 transition-colors">{card.title}</h3>
                <p className="text-sm text-gray-500">{card.desc}</p>
              </Link>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}
