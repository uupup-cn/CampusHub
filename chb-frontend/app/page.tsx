'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { motion, useReducedMotion } from 'framer-motion';
import {
  ArrowRight, MessagesSquare, Coins, Plug, Sparkles,
  Store, ShoppingBag, Wallet, Layers, ShieldCheck, Activity,
} from 'lucide-react';
import Reveal from '@/components/Reveal';
import { apiRequest, formatCHB } from '@/lib/api';

interface PoolData {
  public_pool: { total_supply: number; balance: number };
  official_pool: { balance: number };
}

const marqueeItems = [
  '发帖得 CHB',
  '回复得 CHB',
  '签到得 CHB',
  '被点赞得 CHB',
  '积分集市流转',
  'OAuth2 应用生态',
  '官方积分池收取手续费',
  '总量恒定 5 亿 CHB',
];

const ease = [0.22, 1, 0.36, 1] as const;

export default function Home() {
  const reduce = useReducedMotion();
  const [pools, setPools] = useState<PoolData | null>(null);
  const [liveStats, setLiveStats] = useState<{ items: number; apps: number }>({ items: 0, apps: 0 });

  useEffect(() => {
    apiRequest<PoolData>('/api/chb/pools')
      .then((res) => { if (res.code === 0) setPools(res.data); })
      .catch(() => {});
    apiRequest<{ total_marketplace_orders: number }>('/api/admin/stats')
      .then((res) => { if (res.code === 0) setLiveStats((s) => ({ ...s, items: res.data.total_marketplace_orders })); })
      .catch(() => {});
  }, []);

  const waterLevel = pools ? pools.public_pool.balance / pools.public_pool.total_supply : 0.96;

  const heroWords = ['让知识', '流动', '让价值', '循环'];

  return (
    <div className="overflow-x-hidden">
      {/* Hero */}
      <section className="relative min-h-[calc(100vh-4rem)] flex flex-col justify-center px-6">
        {/* Floating orbs */}
        {!reduce && (
          <>
            <motion.div
              className="absolute w-[42rem] h-[42rem] rounded-full pointer-events-none"
              style={{ background: 'radial-gradient(circle, rgba(139,92,246,0.14), transparent 60%)', top: '-12rem', right: '-10rem' }}
              animate={{ y: [0, -30, 0], x: [0, 18, 0] }}
              transition={{ duration: 14, repeat: Infinity, ease: 'easeInOut' }}
            />
            <motion.div
              className="absolute w-[36rem] h-[36rem] rounded-full pointer-events-none"
              style={{ background: 'radial-gradient(circle, rgba(34,211,238,0.10), transparent 60%)', bottom: '-8rem', left: '-8rem' }}
              animate={{ y: [0, 26, 0], x: [0, -14, 0] }}
              transition={{ duration: 18, repeat: Infinity, ease: 'easeInOut' }}
            />
          </>
        )}

        <motion.div
          className="relative max-w-6xl mx-auto w-full"
          initial={reduce ? false : { opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8 }}
        >
          <motion.div
            className="flex items-center justify-center gap-2 mb-8"
            initial={reduce ? false : { opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.7, ease }}
          >
            <span className="badge badge-purple">
              <Sparkles size={12} />
              CampusHub Ecosystem
            </span>
          </motion.div>

          <h1 className="display-title !text-[clamp(3rem,9vw,8rem)] mb-8 text-center">
            {heroWords.map((w, i) => (
              <motion.span
                key={w}
                className={i % 2 === 1 ? 'block text-transparent bg-clip-text bg-gradient-to-r from-violet-400 via-fuchsia-400 to-cyan-300' : 'block'}
                initial={reduce ? false : { opacity: 0, y: 60, rotate: i % 2 === 1 ? 3 : -2 }}
                animate={{ opacity: 1, y: 0, rotate: 0 }}
                transition={{ duration: 1.1, delay: 0.15 + i * 0.16, ease }}
              >
                {w}
              </motion.span>
            ))}
          </h1>

          <motion.p
            className="text-gray-400 max-w-2xl text-lg leading-relaxed mb-12 mx-auto text-center"
            initial={reduce ? false : { opacity: 0, y: 24 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.85, ease }}
          >
            一个以知识共享为核心的论坛生态：发帖、回复、签到、被点赞都能赚取 CHB。
            积分在集市中流转，应用通过 OAuth2 接入生态，价值在你创造的每个瞬间流动。
          </motion.p>

          <motion.div
            className="flex flex-wrap justify-center gap-4 mb-16"
            initial={reduce ? false : { opacity: 0, y: 24 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 1.05, ease }}
          >
            <Link href="/marketplace" className="btn-primary !px-8 !py-4 text-base">
              <Store size={18} />
              逛逛积分集市
              <ArrowRight size={16} />
            </Link>
            <Link href="/marketplace/apply" className="btn-ghost !px-8 !py-4 text-base">
              <ShoppingBag size={18} />
              申请成为商家
            </Link>
          </motion.div>

          {/* Live stats strip */}
          <motion.div
            className="grid grid-cols-2 md:grid-cols-4 gap-3 max-w-3xl mx-auto"
            initial={reduce ? false : { opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.9, delay: 1.25, ease }}
          >
            {[
              { label: '积分总量', value: '5 亿 CHB', icon: Coins, tint: 'text-amber-300' },
              { label: '公共池存量', value: pools ? formatCHB(pools.public_pool.balance) : '计算中', icon: Wallet, tint: 'text-violet-300' },
              { label: '官方池收益', value: pools ? formatCHB(pools.official_pool.balance) : '计算中', icon: ShieldCheck, tint: 'text-cyan-300' },
              { label: '生态订单', value: String(liveStats.items || 0), icon: Activity, tint: 'text-emerald-300' },
            ].map((s) => {
              const Icon = s.icon;
              return (
                <div key={s.label} className="glass-card !rounded-2xl px-5 py-4">
                  <div className="flex items-center gap-2 mb-1.5">
                    <Icon size={14} className={s.tint} />
                    <span className="text-xs text-gray-500">{s.label}</span>
                  </div>
                  <div className="text-lg font-bold tracking-tight">{s.value}</div>
                </div>
              );
            })}
          </motion.div>
        </motion.div>
      </section>

      {/* Marquee */}
      <section className="border-y border-white/5 py-6 overflow-hidden relative">
        <div className="flex w-max marquee-track">
          {[...marqueeItems, ...marqueeItems].map((item, i) => (
            <span key={i} className="flex items-center gap-8 px-8 text-sm text-gray-500 whitespace-nowrap">
              {item}
              <Sparkles size={12} className="text-violet-500/60" />
            </span>
          ))}
        </div>
      </section>

      {/* Ecosystem */}
      <section className="max-w-7xl mx-auto px-6 py-28">
        <Reveal>
          <div className="flex items-center gap-2 mb-4">
            <Layers size={13} />
            <span className="section-label">Three Pillars</span>
          </div>
          <h2 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-16">
            一个论坛，一套生态
            <span className="block text-2xl md:text-3xl font-semibold text-gray-500 mt-4">知识创造价值，价值滋养生态</span>
          </h2>
        </Reveal>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            {
              icon: MessagesSquare,
              title: '论坛 · 知识共享',
              desc: '发帖、回复、签到、被点赞赚取 CHB。无积分付费设计，没有悬赏提问，让分享回归纯粹。',
              tint: 'from-violet-500/25 to-violet-500/5',
              accent: 'text-violet-300',
            },
            {
              icon: Coins,
              title: '积分 · CHB 经济',
              desc: '总量恒定 5 亿，公共池与官方池逻辑隔绝。集市交易收取 10% 手续费，官方池收益反哺生态治理。',
              tint: 'from-amber-500/25 to-amber-500/5',
              accent: 'text-amber-300',
            },
            {
              icon: Plug,
              title: '应用 · OAuth2 接入',
              desc: '内部应用授权登录论坛账号，实时回调扣减积分。等级门槛可控，一份积分永远不会被重复使用。',
              tint: 'from-cyan-500/25 to-cyan-500/5',
              accent: 'text-cyan-300',
            },
          ].map((card, i) => {
            const Icon = card.icon;
            return (
              <Reveal key={card.title} delay={i * 0.12}>
                <motion.div
                  className="glass-card p-8 h-full group relative overflow-hidden"
                  whileHover={{ y: -8 }}
                  transition={{ type: 'spring', stiffness: 260, damping: 22 }}
                >
                  <div className={`absolute inset-0 bg-gradient-to-br ${card.tint} opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none`} />
                  <div className="relative">
                    <div className="w-12 h-12 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                      <Icon size={22} className={card.accent} />
                    </div>
                    <h3 className="text-xl font-bold mb-3">{card.title}</h3>
                    <p className="text-sm text-gray-400 leading-relaxed">{card.desc}</p>
                  </div>
                </motion.div>
              </Reveal>
            );
          })}
        </div>
      </section>

      {/* Pool visual */}
      <section className="max-w-7xl mx-auto px-6 pb-28">
        <Reveal>
          <div className="gradient-border p-8 md:p-14">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              <div>
                <span className="section-label mb-4 block">Economic Engine</span>
                <h2 className="text-3xl md:text-5xl font-extrabold tracking-tight mb-6">
                  积分池水位
                  <span className="block text-transparent bg-clip-text bg-gradient-to-r from-violet-400 to-cyan-300 mt-2">
                    实时可见的经济脉搏
                  </span>
                </h2>
                <p className="text-gray-400 leading-relaxed mb-8">
                  CHB 总量恒定 5 亿永不增发。公共积分池通过发帖、回复、签到等行为释放给用户，
                  官方积分池仅从集市与第三方应用手续费中获得收益，两池逻辑隔绝、对账可审计。
                </p>
                <div className="flex flex-wrap gap-3">
                  <span className="badge badge-purple">总量 5,000,000,000 CHB</span>
                  <span className="badge badge-cyan">手续费默认 10%</span>
                  <span className="badge badge-green">每日 06:00 重置上限</span>
                </div>
              </div>

              <div>
                <div className="flex items-end justify-between mb-3">
                  <span className="text-sm text-gray-400">公共池水位</span>
                  <span className="text-2xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-violet-400 to-cyan-300">
                    {(waterLevel * 100).toFixed(1)}%
                  </span>
                </div>
                <div className="h-4 rounded-full bg-white/5 overflow-hidden mb-10">
                  <motion.div
                    className="h-full rounded-full bg-gradient-to-r from-violet-500 via-fuchsia-400 to-cyan-300"
                    initial={reduce ? false : { width: 0 }}
                    whileInView={{ width: `${Math.min(100, waterLevel * 100)}%` }}
                    viewport={{ once: true }}
                    transition={{ duration: 1.6, ease }}
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="glass-card p-5">
                    <div className="text-xs text-gray-500 mb-1">官方池</div>
                    <div className="font-bold text-amber-300 text-lg">
                      {pools ? formatCHB(pools.official_pool.balance) : '加载中'}
                    </div>
                  </div>
                  <div className="glass-card p-5">
                    <div className="text-xs text-gray-500 mb-1">生态状态</div>
                    <div className="flex items-center gap-2 font-bold text-emerald-300 text-lg">
                      <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse-glow" />
                      运行中
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Reveal>
      </section>

      {/* CTA */}
      <section className="relative px-6 pb-32">
        <Reveal>
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-6">
              准备好让价值
              <span className="block text-transparent bg-clip-text bg-gradient-to-r from-violet-400 to-cyan-300">流动起来了吗</span>
            </h2>
            <p className="text-gray-400 max-w-xl mx-auto mb-10">
              无论你是内容创作者、商家，还是希望接入生态的应用开发者，
              CampusHub 都为你准备好了完整的积分经济基础设施。
            </p>
            <div className="flex flex-wrap justify-center gap-4">
              <Link href="/marketplace" className="btn-primary !px-8 !py-4">
                <Store size={18} />
                进入集市
              </Link>
              <Link href="/marketplace/apply" className="btn-ghost !px-8 !py-4">
                <ShoppingBag size={18} />
                申请入驻
              </Link>
            </div>
          </div>
        </Reveal>
      </section>
    </div>
  );
}
