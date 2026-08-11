'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { ArrowLeft, ShoppingCart, Package, AlertTriangle, User, Clock, Coins } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, genIdempotencyKey } from '@/lib/api';
import { useToast } from '@/components/Toast';

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

export default function ItemDetailPage() {
  const params = useParams();
  const router = useRouter();
  const toast = useToast();
  const [item, setItem] = useState<Item | null>(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [ordering, setOrdering] = useState(false);
  const [balance, setBalance] = useState(0);

  const itemId = Number(params.id);

  useEffect(() => {
    Promise.all([
      apiRequest<Item>(`/api/marketplace/items/${itemId}`),
      apiRequest<{ balance: number }>('/api/chb/balance'),
    ])
      .then(([itemRes, balRes]) => {
        if (itemRes.code === 0) setItem(itemRes.data);
        if (balRes.code === 0) setBalance(balRes.data.balance);
      })
      .finally(() => setLoading(false));
  }, [itemId]);

  const handleOrder = async () => {
    if (!item) return;
    setOrdering(true);
    try {
      const res = await apiRequest<{ order_no: string; total_amount: number; status: string }>(
        '/api/marketplace/orders',
        {
          method: 'POST',
          body: {
            item_id: item.id,
            quantity,
            idempotency_key: genIdempotencyKey('order'),
          },
        }
      );
      if (res.code === 0) {
        toast('success', `下单成功！订单号：${res.data.order_no}`);
        router.push('/marketplace/orders');
      } else {
        toast('error', res.message);
      }
    } catch {
      toast('error', '网络错误，请稍后重试');
    } finally {
      setOrdering(false);
    }
  };

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-6 py-16">
        <div className="glass-card h-[480px] shimmer" />
      </div>
    );
  }

  if (!item) {
    return (
      <div className="max-w-7xl mx-auto px-6 py-16 text-center">
        <AlertTriangle size={40} className="mx-auto mb-4 text-amber-400" />
        <p className="text-gray-400">商品不存在或已下架</p>
        <button className="btn-ghost mt-6" onClick={() => router.push('/marketplace')}>
          <ArrowLeft size={15} />
          返回集市
        </button>
      </div>
    );
  }

  const canBuy = item.status === 'approved' && item.stock > 0;
  const totalPrice = item.price * quantity;
  const enoughBalance = balance >= totalPrice;

  return (
    <div className="max-w-7xl mx-auto px-6 py-16">
      <button
        onClick={() => router.back()}
        className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
      >
        <ArrowLeft size={16} />
        返回
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        {/* Left - image */}
        <div className="glass-card overflow-hidden">
          <div className="product-thumb h-[420px] relative">
            {item.image_url ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img src={item.image_url} alt={item.title} className="w-full h-full object-cover" />
            ) : (
              <div className="flex flex-col items-center gap-3 text-gray-500">
                <Package size={64} className="animate-float" />
                <span className="text-sm">{item.category || '未分类'}</span>
              </div>
            )}
            <div className="absolute top-4 left-4 flex gap-2">
              <span className="badge badge-purple">{item.category || '其他'}</span>
              {item.status === 'pending' && <span className="badge badge-yellow">待审核</span>}
              {item.status === 'rejected' && <span className="badge badge-red">已拒绝</span>}
            </div>
          </div>
        </div>

        {/* Right - info */}
        <div>
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">{item.title}</h1>

          <div className="flex items-center gap-6 mb-6 text-sm text-gray-500">
            <span className="flex items-center gap-1.5">
              <User size={14} />
              卖家 #{item.seller_id}
            </span>
            <span className="flex items-center gap-1.5">
              <Clock size={14} />
              {formatDate(item.created_at)}
            </span>
            <span className={item.stock > 0 ? 'text-emerald-400' : 'text-red-400'}>
              {item.stock > 0 ? `库存 ${item.stock}` : '已售罄'}
            </span>
          </div>

          <div className="glass-card p-6 mb-8">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-500">单价</span>
              <span className="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-violet-400 to-cyan-300">
                {formatCHB(item.price)}
              </span>
            </div>
            <div className="flex items-center justify-between mb-4">
              <span className="text-sm text-gray-500">我的余额</span>
              <span className={`font-semibold ${enoughBalance ? 'text-emerald-400' : 'text-red-400'}`}>
                {formatCHB(balance)}
              </span>
            </div>

            {/* Quantity */}
            <div className="flex items-center justify-between mb-6">
              <span className="text-sm text-gray-500">购买数量</span>
              <div className="flex items-center gap-3">
                <button
                  className="w-8 h-8 rounded-full border border-white/10 text-gray-400 hover:text-white hover:border-white/30 transition-all disabled:opacity-30"
                  disabled={quantity <= 1}
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                >
                  -
                </button>
                <span className="w-8 text-center font-semibold">{quantity}</span>
                <button
                  className="w-8 h-8 rounded-full border border-white/10 text-gray-400 hover:text-white hover:border-white/30 transition-all disabled:opacity-30"
                  disabled={quantity >= item.stock}
                  onClick={() => setQuantity(Math.min(item.stock, quantity + 1))}
                >
                  +
                </button>
              </div>
            </div>

            <div className="flex items-center justify-between pt-4 border-t border-white/5">
              <span className="text-sm text-gray-500">应付总额</span>
              <span className="text-xl font-bold text-white">{formatCHB(totalPrice)}</span>
            </div>
          </div>

          {!enoughBalance && canBuy && (
            <div className="flex items-center gap-2 text-sm text-amber-400 mb-4">
              <AlertTriangle size={15} />
              余额不足，请先通过论坛活跃赚取 CHB
            </div>
          )}

          <div className="flex gap-4">
            <button
              className="btn-primary flex-1 !py-4 text-base"
              disabled={!canBuy || !enoughBalance || ordering}
              onClick={handleOrder}
            >
              <ShoppingCart size={18} />
              {!canBuy
                ? item.status === 'pending' ? '待审核' : '已售罄'
                : ordering ? '下单中...' : '立即购买'}
            </button>
          </div>

          {item.description && (
            <div className="mt-10">
              <h2 className="section-label mb-4 flex items-center gap-2">
                <Coins size={13} />
                商品描述
              </h2>
              <p className="text-gray-400 leading-relaxed whitespace-pre-line">{item.description}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}