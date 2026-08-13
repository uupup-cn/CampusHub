"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import Link from "next/link";
import { Wallet, TrendingUp, TrendingDown, AlertCircle, MessageCircle, ArrowRight, Package } from "lucide-react";
import { apiRequest, formatCHB } from "@/lib/api";

interface Summary {
  income_7d: number;
  expense_7d: number;
  pending_disputes: number;
  my_disputes: number;
}

export default function DashboardHome() {
  const [summary, setSummary] = useState<Summary | null>(null);

  useEffect(() => {
    apiRequest<Summary>("/api/chb/me/summary")
      .then(res => { if (res.code === 0) setSummary(res.data); })
      .catch(() => {});
  }, []);

  const cards = [
    { label: "7天收入", value: summary ? formatCHB(summary.income_7d) : "加载中", icon: TrendingUp, tint: "text-emerald-300", href: "/dashboard/points" },
    { label: "7天支出", value: summary ? formatCHB(summary.expense_7d) : "加载中", icon: TrendingDown, tint: "text-red-300", href: "/dashboard/points" },
    { label: "待处理争议", value: String(summary?.pending_disputes ?? 0), icon: AlertCircle, tint: "text-amber-300", href: "/dashboard/marketplace/disputes?role=seller" },
    { label: "我发起的争议", value: String(summary?.my_disputes ?? 0), icon: MessageCircle, tint: "text-violet-300", href: "/dashboard/marketplace/disputes?role=buyer" },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">个人中心</h1>

      <div className="grid grid-cols-2 gap-4 mb-8">
        {cards.map((card, i) => {
          const Icon = card.icon;
          return (
            <motion.div key={card.label} initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: i * 0.08 }}>
              <Link href={card.href} className="glass-card p-5 block hover:border-white/20 transition-colors">
                <div className="flex items-center gap-2 mb-2">
                  <Icon size={15} className={card.tint} />
                  <span className="text-xs text-gray-500">{card.label}</span>
                </div>
                <div className="text-lg font-bold">{card.value}</div>
              </Link>
            </motion.div>
          );
        })}
      </div>

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
