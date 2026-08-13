'use client';

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Wallet, TrendingUp, Package, Activity, ArrowRight } from 'lucide-react';
import Link from 'next/link';
import { apiRequest, formatCHB } from '@/lib/api';

interface SummaryData {
  balance: number;
  today_earned: number;
  daily_cap: number;
  orders_count: number;
}

interface Transaction {
  id: number;
  type: string;
  amount: number;
  description: string;
  created_at: string;
}

export default function DashboardHome() {
  const [summary, setSummary] = useState<SummaryData | null>(null);
  const [recentTx, setRecentTx] = useState<Transaction[]>([]);

  useEffect(() => {
    apiRequest<{ balance: number }>('/api/chb/balance')
      .then(res => { if (res.code === 0) setSummary(prev => ({ balance: res.data.balance, today_earned: prev?.today_earned ?? 0, daily_cap: prev?.daily_cap ?? 0, orders_count: prev?.orders_count ?? 0 })); })
      .catch(() => {});

    apiRequest<{ total_marketplace_orders: number }>('/api/admin/stats')
      .then(res => { if (res.code === 0) setSummary(prev => ({ balance: prev?.balance ?? 0, today_earned: prev?.today_earned ?? 0, daily_cap: prev?.daily_cap ?? 0, orders_count: res.data.total_marketplace_orders })); })
      .catch(() => {});
  }, []);

  const cards = [
    { label: '积分余额', value: summary ? formatCHB(summary.balance) : '加载中', icon: Wallet, tint: 'text-amber-300' },
    { label: '今日已获', value: '0 CHB', icon: TrendingUp, tint: 'text-emerald-300' },
    { label: '集市订单', value: String(summary?.orders_count ?? 0), icon: Package, tint: 'text-violet-300' },
    { label: '授权应用', value: '0', icon: Activity, tint: 'text-cyan-300' },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">个人中心</h1>

      {/* Stats cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        {cards.map((card, i) => {
          const Icon = card.icon;
          return (
            <motion.div
              key={card.label}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.08 }}
              className="glass-card p-5"
            >
              <div className="flex items-center gap-2 mb-2">
                <Icon size={15} className={card.tint} />
                <span className="text-xs text-gray-500">{card.label}</span>
              </div>
              <div className="text-lg font-bold">{card.value}</div>
            </motion.div>
          );
        })}
      </div>

      {/* Quick links */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Link href="/dashboard/points" className="glass-card p-6 group hover:border-white/20 transition-colors">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium">积分明细</span>
            <ArrowRight size={16} className="text-gray-600 group-hover:text-violet-400 transition-colors" />
          </div>
          <p className="text-xs text-gray-500">查看全部积分流水记录</p>
        </Link>
        <Link href="/dashboard/marketplace" className="glass-card p-6 group hover:border-white/20 transition-colors">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium">集市管理</span>
            <ArrowRight size={16} className="text-gray-600 group-hover:text-violet-400 transition-colors" />
          </div>
          <p className="text-xs text-gray-500">管理商品、查看订单</p>
        </Link>
      </div>
    </div>
  );
}
