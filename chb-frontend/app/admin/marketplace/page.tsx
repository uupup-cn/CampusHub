'use client';

import { useCallback, useEffect, useState } from 'react';
import {
  PackageCheck, Store, Check, X, Clock, RefreshCcw,
} from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';
import { useToast } from '@/components/Toast';

interface MerchantApp {
  id: number;
  discourse_user_id: number;
  shop_name: string;
  description: string | null;
  status: string;
  created_at: string;
}

interface Item {
  id: number;
  seller_id: number;
  title: string;
  price: number;
  stock: number;
  status: string;
  category: string | null;
  created_at: string;
}

type Tab = 'items' | 'applications';

export default function AdminMarketplacePage() {
  const [tab, setTab] = useState<Tab>('items');
  const [items, setItems] = useState<Item[]>([]);
  const [apps, setApps] = useState<MerchantApp[]>([]);
  const [loading, setLoading] = useState(true);
  const toast = useToast();

  const fetchData = useCallback(async () => {
    try {
      if (tab === 'items') {
        const res = await apiRequest<Paginated<Item>>('/api/admin/marketplace/items?page=1&page_size=50');
        if (res.code === 0) setItems(res.data.items);
      } else {
        const res = await apiRequest<Paginated<MerchantApp>>('/api/admin/marketplace/applications?page=1&page_size=50');
        if (res.code === 0) setApps(res.data.items);
      }
    } finally {
      setLoading(false);
    }
  }, [tab]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const reviewItem = async (id: number, status: 'approved' | 'rejected') => {
    const res = await apiRequest(`/api/admin/marketplace/items/${id}`, {
      method: 'PUT',
      body: { status, review_comment: status === 'approved' ? '审核通过' : '审核拒绝' },
    });
    if (res.code === 0) {
      toast('success', status === 'approved' ? '商品已通过' : '商品已拒绝');
      fetchData();
    } else {
      toast('error', res.message);
    }
  };

  const reviewApp = async (id: number, status: 'approved' | 'rejected') => {
    const res = await apiRequest(`/api/admin/marketplace/applications/${id}`, {
      method: 'PUT',
      body: { status, review_comment: status === 'approved' ? '审核通过' : '审核拒绝' },
    });
    if (res.code === 0) {
      toast('success', status === 'approved' ? '入驻申请已通过' : '入驻申请已拒绝');
      fetchData();
    } else {
      toast('error', res.message);
    }
  };

  return (
    <div className="max-w-6xl mx-auto px-6 py-16">
      <div className="flex items-center justify-between mb-10">
        <div>
          <span className="section-label flex items-center gap-2 mb-3">
            <PackageCheck size={13} />
            Marketplace Review
          </span>
          <h1 className="text-3xl font-bold tracking-tight">集市审核</h1>
        </div>
        <button className="btn-ghost !px-4 !py-2" onClick={fetchData}>
          <RefreshCcw size={15} />
          刷新
        </button>
      </div>

      <div className="flex gap-2 mb-8">
        <button
          className={`px-5 py-2.5 rounded-full text-sm transition-all ${
            tab === 'items' ? 'bg-violet-500/20 text-violet-300 border border-violet-500/30' : 'text-gray-400 border border-white/5 hover:border-white/20'
          }`}
          onClick={() => setTab('items')}
        >
          商品审核
        </button>
        <button
          className={`px-5 py-2.5 rounded-full text-sm transition-all ${
            tab === 'applications' ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/30' : 'text-gray-400 border border-white/5 hover:border-white/20'
          }`}
          onClick={() => setTab('applications')}
        >
          入驻申请
        </button>
      </div>

      {loading ? (
        <div className="space-y-4">{['','',''].map((_,i) => <div key={i} className="glass-card h-20 shimmer" />)}</div>
      ) : tab === 'items' ? (
        items.length === 0 ? (
          <div className="glass-card p-16 text-center">
            <PackageCheck size={40} className="mx-auto mb-4 text-gray-600" />
            <p className="text-gray-500">暂无待审核商品</p>
          </div>
        ) : (
          <div className="space-y-4">
            {items.map((item) => (
              <div key={item.id} className="glass-card p-5 flex items-center gap-5">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 mb-1">
                    <h3 className="font-semibold truncate">{item.title}</h3>
                    <span className="badge badge-yellow"><Clock size={12} />待审核</span>
                  </div>
                  <p className="text-sm text-gray-500">
                    卖家 #{item.seller_id} · {item.category || '未分类'} · 库存 {item.stock} · {formatDate(item.created_at)}
                  </p>
                </div>
                <div className="font-bold text-violet-300">{formatCHB(item.price)}</div>
                <div className="flex gap-2">
                  <button
                    className="w-9 h-9 rounded-full bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center hover:bg-emerald-500/20 transition-all"
                    onClick={() => reviewItem(item.id, 'approved')}
                    title="通过"
                  >
                    <Check size={15} className="text-emerald-400" />
                  </button>
                  <button
                    className="w-9 h-9 rounded-full bg-red-500/10 border border-red-500/30 flex items-center justify-center hover:bg-red-500/20 transition-all"
                    onClick={() => reviewItem(item.id, 'rejected')}
                    title="拒绝"
                  >
                    <X size={15} className="text-red-400" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )
      ) : apps.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Store size={40} className="mx-auto mb-4 text-gray-600" />
          <p className="text-gray-500">暂无入驻申请</p>
        </div>
      ) : (
        <div className="space-y-4">
          {apps.map((app) => (
            <div key={app.id} className="glass-card p-5 flex items-center gap-5">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-3 mb-1">
                  <h3 className="font-semibold">{app.shop_name}</h3>
                  <span className="badge badge-yellow"><Clock size={12} />待审核</span>
                </div>
                <p className="text-sm text-gray-500 truncate">
                  用户 #{app.discourse_user_id} · {app.description || '无介绍'} · {formatDate(app.created_at)}
                </p>
              </div>
              <div className="flex gap-2">
                <button
                  className="w-9 h-9 rounded-full bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center hover:bg-emerald-500/20 transition-all"
                  onClick={() => reviewApp(app.id, 'approved')}
                  title="通过"
                >
                  <Check size={15} className="text-emerald-400" />
                </button>
                <button
                  className="w-9 h-9 rounded-full bg-red-500/10 border border-red-500/30 flex items-center justify-center hover:bg-red-500/20 transition-all"
                  onClick={() => reviewApp(app.id, 'rejected')}
                  title="拒绝"
                >
                  <X size={15} className="text-red-400" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
