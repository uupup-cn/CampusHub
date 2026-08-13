"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { AlertCircle, Check, X } from "lucide-react";
import { apiRequest, formatCHB, formatDate, type Paginated } from "@/lib/api";

interface Dispute {
  id: number;
  order_id: number;
  buyer_id: number;
  seller_id: number;
  status: string;
  buyer_reason: string;
  refund_amount: number;
  created_at: string;
  seller_action: string | null;
  seller_reason: string | null;
}

const statusLabels: Record<string, string> = {
  pending: "等待卖家处理",
  accepted: "已退款",
  rejected: "等待管理员审核",
  auto_refunded: "超时自动退款",
  admin_buyer_win: "管理员判定退款",
  admin_seller_win: "管理员判定卖家胜",
  closed: "已关闭",
};

function DisputesContent() {
  const searchParams = useSearchParams();
  const [disputes, setDisputes] = useState<Dispute[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const role = searchParams.get("role") || "buyer";

  const fetchDisputes = useCallback(async () => {
    try {
      const res = await apiRequest<Paginated<Dispute>>(`/api/marketplace/disputes?role=${role}&page=${page}&page_size=20`);
      if (res.code === 0) {
        setDisputes(res.data.items);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  }, [role, page]);

  useEffect(() => { fetchDisputes(); }, [fetchDisputes]);

  const handleAccept = async (id: number) => {
    await apiRequest(`/api/marketplace/disputes/${id}/accept`, { method: "PUT" });
    fetchDisputes();
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">争议管理</h1>
      <div className="flex gap-4 mb-6">
        <Link href="/dashboard/marketplace" className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white transition-colors">我的商品</Link>
        <Link href="/dashboard/marketplace/orders" className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white transition-colors">我的订单</Link>
        <span className="px-4 py-2 rounded-lg text-sm bg-violet-500/15 text-violet-300 font-medium">争议管理</span>
      </div>
      {loading ? (
        <div className="space-y-3">{[1, 2, 3].map(i => <div key={i} className="glass-card h-20 shimmer" />)}</div>
      ) : disputes.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <AlertCircle size={32} className="mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500">暂无争议</p>
        </div>
      ) : (
        <div className="space-y-3">
          {disputes.map(d => (
            <div key={d.id} className="glass-card p-5">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <span className="text-sm font-medium">争议 #{d.id}</span>
                  <span className="badge !text-[10px] badge-gray">{statusLabels[d.status] || d.status}</span>
                </div>
                <span className="text-sm font-bold text-amber-300">{formatCHB(d.refund_amount)}</span>
              </div>
              <p className="text-xs text-gray-500 mb-3">原因: {d.buyer_reason}</p>
              <p className="text-xs text-gray-600">{formatDate(d.created_at)}</p>
              {role === "seller" && d.status === "pending" && (
                <div className="flex gap-3 mt-4">
                  <button onClick={() => handleAccept(d.id)} className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-500/15 text-emerald-300 text-xs hover:bg-emerald-500/25 transition-all">
                    <Check size={13} /> 同意退款
                  </button>
                  <Link href={`/dashboard/marketplace/disputes/${d.id}?action=reject`} className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-red-500/15 text-red-300 text-xs hover:bg-red-500/25 transition-all">
                    <X size={13} /> 拒绝退款
                  </Link>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default function DisputesPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center"><div className="glass-card h-20 w-64 shimmer" /></div>}>
      <DisputesContent />
    </Suspense>
  );
}
