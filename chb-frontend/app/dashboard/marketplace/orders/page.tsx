"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { ShoppingBag, AlertCircle, X } from "lucide-react";
import { apiRequest, formatCHB, formatDate, type Paginated } from "@/lib/api";

interface Order {
  id: number;
  item_title: string;
  quantity: number;
  total_amount: number;
  status: string;
  created_at: string;
  dispute_status: string | null;
}

export default function DashboardOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<"buyer" | "seller">("buyer");
  const [page, setPage] = useState(1);
  const [disputeModal, setDisputeModal] = useState<{ orderId: number; orderTitle: string } | null>(null);
  const [disputeReason, setDisputeReason] = useState("");
  const [disputeImages, setDisputeImages] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const fetchOrders = useCallback(async () => {
    try {
      const res = await apiRequest<Paginated<Order>>(`/api/marketplace/orders?role=${tab}&page=${page}&page_size=10`);
      if (res.code === 0) {
        setOrders(res.data.items as Order[]);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  }, [tab, page]);

  useEffect(() => { fetchOrders(); }, [fetchOrders]);

  const canDispute = (order: Order) => {
    if (tab !== "buyer") return false;
    if (order.dispute_status) return false;
    const created = new Date(order.created_at).getTime();
    return Date.now() - created < 3 * 24 * 60 * 60 * 1000;
  };

  const submitDispute = async () => {
    if (!disputeModal || !disputeReason.trim()) return;
    setSubmitting(true);
    try {
      const images = disputeImages ? disputeImages.split("\n").map(s => s.trim()).filter(Boolean) : [];
      const res = await apiRequest(`/api/marketplace/orders/${disputeModal.orderId}/dispute`, {
        method: "POST",
        body: { reason: disputeReason, images },
      });
      if (res.code === 0) {
        setDisputeModal(null);
        setDisputeReason("");
        setDisputeImages("");
        fetchOrders();
      }
    } finally {
      setSubmitting(false);
    }
  };

  const totalPages = Math.max(1, Math.ceil(total / 10));

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold">集市管理</h1>
      </div>

      <div className="flex gap-4 mb-6">
        <Link href="/dashboard/marketplace" className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white transition-colors">我的商品</Link>
        <span className="px-4 py-2 rounded-lg text-sm bg-violet-500/15 text-violet-300 font-medium">我的订单</span>
        <Link href="/dashboard/marketplace/disputes" className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white transition-colors">争议管理</Link>
      </div>

      <div className="flex gap-2 mb-6">
        <button onClick={() => { setTab("buyer"); setPage(1); }} className={"px-3 py-1.5 rounded-full text-xs transition-all " + (tab === "buyer" ? "bg-violet-500/20 text-violet-300 border border-violet-500/30" : "text-gray-400 border border-white/5")}>我买入的</button>
        <button onClick={() => { setTab("seller"); setPage(1); }} className={"px-3 py-1.5 rounded-full text-xs transition-all " + (tab === "seller" ? "bg-violet-500/20 text-violet-300 border border-violet-500/30" : "text-gray-400 border border-white/5")}>我卖出的</button>
      </div>

      {loading ? (
        <div className="space-y-3">{[1, 2, 3].map(i => <div key={i} className="glass-card h-16 shimmer" />)}</div>
      ) : orders.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <ShoppingBag size={32} className="mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500">暂无订单</p>
        </div>
      ) : (
        <div className="space-y-3">
          {orders.map(order => (
            <div key={order.id} className="glass-card p-4 flex items-center justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-3 mb-1">
                  <span className="font-medium truncate">{order.item_title}</span>
                  <span className="badge !text-[10px] badge-gray">{order.status}</span>
                  {order.dispute_status && <span className="badge !text-[10px] badge-yellow">争议中</span>}
                </div>
                <div className="text-xs text-gray-500">数量 {order.quantity} / {formatDate(order.created_at)}</div>
              </div>
              <div className="flex items-center gap-3">
                <div className="text-right">
                  <div className="text-sm font-bold text-amber-300">{formatCHB(order.total_amount)}</div>
                </div>
                {canDispute(order) && (
                  <button onClick={() => setDisputeModal({ orderId: order.id, orderTitle: order.item_title })} className="flex items-center gap-1.5 px-3 py-1.5 rounded-full border border-amber-500/30 text-xs text-amber-300 hover:bg-amber-500/10 transition-all">
                    <AlertCircle size={13} />
                    发起争议
                  </button>
                )}
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

      {/* 发起争议弹窗 */}
      {disputeModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" onClick={() => setDisputeModal(null)}>
          <div className="glass-card p-8 w-full max-w-md mx-4" onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-bold">发起订单争议</h2>
              <button onClick={() => setDisputeModal(null)} className="text-gray-500 hover:text-white"><X size={18} /></button>
            </div>
            <p className="text-xs text-gray-500 mb-4">订单: {disputeModal.orderTitle}</p>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">退款原因</label>
              <textarea className="input" placeholder="请描述退款原因..." value={disputeReason} onChange={e => setDisputeReason(e.target.value)} rows={4} />
            </div>
            <div className="mb-6">
              <label className="block text-sm font-medium mb-2">图片URL（非必填，每行一个）</label>
              <textarea className="input" placeholder="https://example.com/image1.png" value={disputeImages} onChange={e => setDisputeImages(e.target.value)} rows={2} />
            </div>
            <button className="btn-primary w-full" onClick={submitDispute} disabled={submitting || !disputeReason.trim()}>
              {submitting ? "提交中..." : "提交争议"}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
