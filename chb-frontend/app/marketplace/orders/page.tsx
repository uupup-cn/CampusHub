'use client';

import { useEffect, useState } from 'react';
import { ArrowLeft, ShoppingBag, Package, CheckCircle2, Clock, XCircle, RefreshCcw } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';

interface Order {
  id: number;
  order_no: string;
  item_id: number;
  buyer_id: number;
  seller_id: number;
  quantity: number;
  unit_price: number;
  total_amount: number;
  fee: number;
  net_amount: number;
  status: string;
  created_at: string;
}

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [role, setRole] = useState<'buyer' | 'seller'>('buyer');

  useEffect(() => {
    fetchOrders();
  }, [role]);

  async function fetchOrders() {
    setLoading(true);
    try {
      const res = await apiRequest<Paginated<Order>>(`/api/marketplace/orders?role=${role}&page=1&page_size=50`);
      if (res.code === 0) {
        setOrders(res.data.items);
      }
    } finally {
      setLoading(false);
    }
  }

  const statusBadge = (status: string) => {
    switch (status) {
      case 'completed': return <span className="badge badge-green"><CheckCircle2 size={12} />已完成</span>;
      case 'pending': return <span className="badge badge-yellow"><Clock size={12} />待处理</span>;
      case 'cancelled': return <span className="badge badge-red"><XCircle size={12} />已取消</span>;
      default: return <span className="badge badge-gray">{status}</span>;
    }
  };

  return (
    <div className="max-w-4xl mx-auto px-6 py-16">
      <button
        onClick={() => window.history.back()}
        className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-10"
      >
        <ArrowLeft size={16} />
        返回
      </button>

      <div className="flex items-center justify-between mb-10">
        <div>
          <span className="section-label flex items-center gap-2 mb-3">
            <ShoppingBag size={13} />
            My Orders
          </span>
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">我的订单</h1>
        </div>
        <div className="flex gap-2">
          <button
            className={`px-4 py-2 rounded-full text-sm transition-all ${
              role === 'buyer' ? 'bg-violet-500/20 text-violet-300 border border-violet-500/30' : 'text-gray-400 border border-white/5 hover:border-white/20'
            }`}
            onClick={() => setRole('buyer')}
          >
            作为买家
          </button>
          <button
            className={`px-4 py-2 rounded-full text-sm transition-all ${
              role === 'seller' ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/30' : 'text-gray-400 border border-white/5 hover:border-white/20'
            }`}
            onClick={() => setRole('seller')}
          >
            作为卖家
          </button>
        </div>
      </div>

      {loading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => <div key={i} className="glass-card h-24 shimmer" />)}
        </div>
      ) : orders.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Package size={40} className="mx-auto mb-4 text-gray-600" />
          <p className="text-gray-500">暂无{role === 'buyer' ? '购买' : '卖出'}订单</p>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div key={order.id} className="glass-card p-5">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <span className="font-mono text-sm text-gray-400">{order.order_no}</span>
                  {statusBadge(order.status)}
                </div>
                <span className="text-xs text-gray-600">{formatDate(order.created_at)}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="text-sm text-gray-400">
                  商品 #{order.item_id} · 数量 x{order.quantity} · 单价 {formatCHB(order.unit_price)}
                </div>
                <div className="text-right">
                  <div className="font-bold">{formatCHB(order.total_amount)}</div>
                  <div className="text-xs text-gray-600">手续费 {formatCHB(order.fee)}</div>
                </div>
              </div>
            </div>
          ))}
          <button
            onClick={fetchOrders}
            className="w-full py-3 rounded-xl border border-white/5 text-sm text-gray-500 hover:text-white hover:border-white/20 transition-all flex items-center justify-center gap-2"
          >
            <RefreshCcw size={14} />
            刷新
          </button>
        </div>
      )}
    </div>
  );
}