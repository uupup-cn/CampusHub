'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { ArrowRight, Package, PlusCircle, Clock, CheckCircle2, XCircle, Pencil } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';

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

export default function MyItemsPage() {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiRequest<Paginated<Item>>('/api/marketplace/items/mine?page=1&page_size=50')
      .then((res) => {
        if (res.code === 0) {
          setItems(res.data.items);
        }
      })
      .finally(() => setLoading(false));
  }, []);

  const statusBadge = (status: string) => {
    switch (status) {
      case 'approved': return <span className="badge badge-green"><CheckCircle2 size={12} />已上架</span>;
      case 'pending': return <span className="badge badge-yellow"><Clock size={12} />待审核</span>;
      case 'rejected': return <span className="badge badge-red"><XCircle size={12} />已拒绝</span>;
      default: return <span className="badge badge-gray">{status}</span>;
    }
  };

  return (
    <div className="max-w-4xl mx-auto px-6 py-16">
      <div className="flex items-center justify-between mb-10">
        <div>
          <span className="section-label flex items-center gap-2 mb-3">
            <Package size={13} />
            My Items
          </span>
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">我的商品</h1>
        </div>
        <Link href="/marketplace/apply" className="btn-primary !py-3">
          <PlusCircle size={16} />
          发布商品
        </Link>
      </div>

      {loading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => <div key={i} className="glass-card h-24 shimmer" />)}
        </div>
      ) : items.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Package size={40} className="mx-auto mb-4 text-gray-600" />
          <p className="text-gray-500 mb-2">还没有发布任何商品</p>
          <Link href="/marketplace/apply" className="text-violet-400 hover:text-violet-300 text-sm inline-flex items-center gap-1.5">
            去发布第一个商品
            <ArrowRight size={14} />
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {items.map((item) => (
            <Link
              key={item.id}
              href={`/marketplace/${item.id}`}
              className="glass-card p-5 flex items-center gap-5 hover:translate-x-1 transition-transform"
            >
              <div className="w-16 h-16 rounded-xl product-thumb shrink-0">
                <Package size={22} className="text-gray-500" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-3 mb-1">
                  <h3 className="font-semibold truncate">{item.title}</h3>
                  {statusBadge(item.status)}
                </div>
                <p className="text-sm text-gray-500 truncate">
                  {item.category || '未分类'} · 库存 {item.stock} · {formatDate(item.created_at)}
                </p>
              </div>
              <div className="text-right">
                <div className="font-bold text-violet-300">{formatCHB(item.price)}</div>
                <div className="text-xs text-gray-600 mt-1 flex items-center gap-1 justify-end">
                  <Pencil size={11} />
                  查看详情
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
