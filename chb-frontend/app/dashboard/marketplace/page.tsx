'use client';

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { Package, PlusCircle, Store } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';

interface Item {
  id: number;
  title: string;
  price: number;
  stock: number;
  status: string;
  category: string | null;
  created_at: string;
}

export default function DashboardMarketplacePage() {
  const [items, setItems] = useState<Item[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  const fetchItems = useCallback(async () => {
    try {
      const res = await apiRequest<Paginated<Item>>(`/api/marketplace/items/mine?page=${page}&page_size=10`);
      if (res.code === 0) {
        setItems(res.data.items);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => { fetchItems(); }, [fetchItems]);

  const totalPages = Math.max(1, Math.ceil(total / 10));

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold">集市管理</h1>
        <Link href="/dashboard/marketplace/apply" className="btn-primary !py-2 !px-4 text-sm">
          <PlusCircle size={15} />
          申请入驻
        </Link>
      </div>

      <div className="flex gap-4 mb-6">
        <Link href="/dashboard/marketplace" className="px-4 py-2 rounded-lg text-sm bg-violet-500/15 text-violet-300 font-medium">我的商品</Link>
        <Link href="/dashboard/marketplace/orders" className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white transition-colors">我的订单</Link>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[1, 2, 3].map(i => <div key={i} className="glass-card h-16 shimmer" />)}
        </div>
      ) : items.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <Package size={32} className="mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500 mb-2">暂无商品</p>
          <p className="text-sm text-gray-600">在集市中发布第一个商品</p>
        </div>
      ) : (
        <div className="space-y-3">
          {items.map(item => (
            <div key={item.id} className="glass-card p-4 flex items-center justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-3 mb-1">
                  <span className="font-medium truncate">{item.title}</span>
                  <span className={'badge !text-[10px] ' + (item.status === 'approved' ? 'badge-green' : 'badge-yellow')}>{item.status}</span>
                </div>
                <div className="text-xs text-gray-500">{item.category || '未分类'} / 库存 {item.stock}</div>
              </div>
              <div className="text-right">
                <div className="text-sm font-bold text-amber-300">{formatCHB(item.price)}</div>
                <div className="text-xs text-gray-600">{formatDate(item.created_at)}</div>
              </div>
            </div>
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-3 mt-8">
          <button className="btn-ghost !px-4 !py-2 text-sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>上一页</button>
          <span className="text-sm text-gray-500">{page} / {totalPages}</span>
          <button className="btn-ghost !px-4 !py-2 text-sm" disabled={page >= totalPages} onClick={() => setPage(page + 1)}>下一页</button>
        </div>
      )}
    </div>
  );
}
