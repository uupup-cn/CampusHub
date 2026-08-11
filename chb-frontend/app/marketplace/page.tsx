'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { Search, Filter, ArrowUpDown, ShoppingBag, PlusCircle, Package, Sparkles, TrendingUp, Layers } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';
import Reveal from '@/components/Reveal';

interface Item {
  id: number;
  seller_id: number;
  title: string;
  description: string | null;
  price: number;
  stock: number;
  status: string;
  category: string | null;
  image_url: string | null;
  created_at: string;
}

const categories = ['全部', '周边', '虚拟', '服务', '学习', '生活'];

export default function MarketplacePage() {
  const [items, setItems] = useState<Item[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [keyword, setKeyword] = useState('');
  const [category, setCategory] = useState('全部');
  const [sort, setSort] = useState('newest');
  const [page, setPage] = useState(1);
  const [stats, setStats] = useState({ items: 0, categories: 0, total_chb: 0 });

  useEffect(() => {
    fetchItems();
  }, [keyword, category, sort, page]);

  async function fetchItems() {
    setLoading(true);
    try {
      const params = new URLSearchParams({ page: String(page), page_size: '12', sort });
      if (keyword) params.set('keyword', keyword);
      if (category !== '全部') params.set('category', category);

      const res = await apiRequest<Paginated<Item>>(`/api/marketplace/items?${params}`);
      if (res.code === 0) {
        setItems(res.data.items);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  }

  // Load pool stats
  useEffect(() => {
    apiRequest<{ public_pool: { total_supply: number; balance: number } }>('/api/chb/pools')
      .then((res) => {
        if (res.code === 0) {
          setStats({
            items: total,
            categories: categories.length - 1,
            total_chb: res.data.public_pool.balance,
          });
        }
      })
      .catch(() => {});
  }, [total]);

  const totalPages = Math.max(1, Math.ceil(total / 12));

  return (
    <div className="max-w-7xl mx-auto px-6 py-16">
      {/* Hero */}
      <section className="relative mb-16">
        <div className="flex items-center gap-2 mb-6">
          <span className="section-label flex items-center gap-2">
            <Layers size={13} />
            CampusHub Marketplace
          </span>
        </div>
        <h1 className="display-title mb-6">
          积分集市
          <span className="block text-2xl md:text-4xl font-semibold tracking-tight bg-none bg-clip-border mt-3"
            style={{ WebkitTextFillColor: 'var(--text-dim)', fontSize: 'clamp(1.1rem, 2.5vw, 1.8rem)', fontWeight: 500, letterSpacing: '-0.01em' }}>
            用 CHB 兑换你的校园好物
          </span>
        </h1>
        <p className="text-gray-400 max-w-2xl leading-relaxed mb-10">
          在这里，每一次发帖、回复、签到获得的积分都能变成实实在在的周边、虚拟服务和知识商品。
          创作者入驻集市，让价值流动起来。
        </p>

        {/* Stats */}
        <Reveal>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-10">
            <div className="glass-card stat-card">
              <div className="flex items-center gap-2 mb-2">
                <ShoppingBag size={15} className="text-violet-400" />
                <span className="text-xs text-gray-500">在售商品</span>
              </div>
              <div className="stat-value">{stats.items}</div>
            </div>
            <div className="glass-card stat-card">
              <div className="flex items-center gap-2 mb-2">
                <Package size={15} className="text-cyan-400" />
                <span className="text-xs text-gray-500">商品分类</span>
              </div>
              <div className="stat-value">{stats.categories}</div>
            </div>
            <div className="glass-card stat-card col-span-2 md:col-span-1">
              <div className="flex items-center gap-2 mb-2">
                <TrendingUp size={15} className="text-amber-400" />
                <span className="text-xs text-gray-500">公共池存量</span>
              </div>
              <div className="stat-value text-amber-300">{formatCHB(stats.total_chb)}</div>
            </div>
          </div>
        </Reveal>

        {/* Actions */}
        <Reveal delay={0.1}>
          <div className="flex flex-wrap gap-4">
            <Link href="/marketplace/apply" className="btn-primary">
              <PlusCircle size={16} />
              申请入驻
            </Link>
            <Link href="/marketplace/my-items" className="btn-ghost">
              <Package size={16} />
              我的商品
            </Link>
            <Link href="/marketplace/orders" className="btn-ghost">
              <ShoppingBag size={16} />
              我的订单
            </Link>
          </div>
        </Reveal>
      </section>

      {/* Filters */}
      <section className="mb-10">
        <div className="flex flex-col md:flex-row md:items-center gap-4 justify-between">
          <div className="flex items-center gap-2 relative max-w-md w-full">
            <Search size={16} className="absolute left-4 text-gray-500" />
            <input
              className="input !pl-11"
              placeholder="搜索商品..."
              value={keyword}
              onChange={(e) => { setKeyword(e.target.value); setPage(1); }}
            />
          </div>
          <div className="flex items-center gap-3 flex-wrap">
            <Filter size={15} className="text-gray-500" />
            {categories.map((c) => (
              <button
                key={c}
                onClick={() => { setCategory(c); setPage(1); }}
                className={`px-3.5 py-1.5 rounded-full text-xs transition-all ${
                  category === c
                    ? 'bg-violet-500/20 text-violet-300 border border-violet-500/30'
                    : 'text-gray-400 border border-white/5 hover:border-white/20'
                }`}
              >
                {c}
              </button>
            ))}
          </div>
          <button
            onClick={() => setSort(sort === 'newest' ? 'price_asc' : sort === 'price_asc' ? 'price_desc' : 'newest')}
            className="flex items-center gap-2 px-4 py-2 rounded-full border border-white/10 text-xs text-gray-400 hover:text-white transition-all"
          >
            <ArrowUpDown size={14} />
            {sort === 'newest' ? '最新发布' : sort === 'price_asc' ? '价格从低到高' : '价格从高到低'}
          </button>
        </div>
      </section>

      {/* Grid */}
      <section>
        {loading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="glass-card h-[320px] shimmer" />
            ))}
          </div>
        ) : items.length === 0 ? (
          <div className="text-center py-24">
            <Sparkles size={40} className="mx-auto mb-4 text-gray-600" />
            <p className="text-gray-500 mb-2">暂无商品</p>
            <p className="text-sm text-gray-600">去发布第一个商品吧</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((item, i) => (
              <motion.div
                key={item.id}
                initial={{ opacity: 0, y: 30 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: '-40px' }}
                transition={{ duration: 0.55, delay: (i % 6) * 0.07, ease: [0.22, 1, 0.36, 1] }}
              >
                <Link
                  href={`/marketplace/${item.id}`}
                  className="glass-card overflow-hidden group block"
                >
                  <div className="product-thumb h-44 relative overflow-hidden">
                    {item.image_url ? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img src={item.image_url} alt={item.title} className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                    ) : (
                      <div className="flex flex-col items-center gap-2 text-gray-500">
                        <ShoppingBag size={32} className="group-hover:scale-110 group-hover:text-violet-400 transition-all duration-300" />
                        <span className="text-xs">{item.category || '未分类'}</span>
                      </div>
                    )}
                    <div className="absolute top-3 right-3">
                      <span className="badge badge-purple">{item.category || '其他'}</span>
                    </div>
                  </div>
                  <div className="p-5">
                    <h3 className="font-semibold text-lg mb-2 line-clamp-1 group-hover:text-violet-300 transition-colors">
                      {item.title}
                    </h3>
                    <p className="text-sm text-gray-500 mb-4 line-clamp-2 min-h-[2.5rem]">
                      {item.description || '暂无描述'}
                    </p>
                    <div className="flex items-center justify-between">
                      <span className="text-lg font-bold text-transparent bg-clip-text bg-gradient-to-r from-violet-400 to-cyan-300">
                        {formatCHB(item.price)}
                      </span>
                      <span className="text-xs text-gray-500">
                        <span className={item.stock > 0 ? 'text-emerald-400' : 'text-red-400'}>
                          {item.stock > 0 ? `库存 ${item.stock}` : '已售罄'}
                        </span>
                        {' · '}
                        {formatDate(item.created_at)}
                      </span>
                    </div>
                  </div>
                </Link>
              </motion.div>
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-3 mt-12">
            <button
              className="btn-ghost !px-4 !py-2 text-sm"
              disabled={page <= 1}
              onClick={() => setPage(page - 1)}
            >
              上一页
            </button>
            <span className="text-sm text-gray-500">
              {page} / {totalPages}
            </span>
            <button
              className="btn-ghost !px-4 !py-2 text-sm"
              disabled={page >= totalPages}
              onClick={() => setPage(page + 1)}
            >
              下一页
            </button>
          </div>
        )}
      </section>
    </div>
  );
}
