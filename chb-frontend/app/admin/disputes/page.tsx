"use client";

import { useCallback, useEffect, useState } from "react";
import { AlertCircle, Gavel } from "lucide-react";
import { apiRequest, formatCHB, formatDate, type Paginated } from "@/lib/api";

interface AdminDisputeItem {
  dispute: {
    id: number;
    order_id: number;
    buyer_id: number;
    seller_id: number;
    status: string;
    buyer_reason: string;
    seller_reason: string | null;
    admin_note: string | null;
    refund_amount: number;
    created_at: string;
    resolved_at: string | null;
  };
  item_title: string;
}

const statusLabels: Record<string, string> = {
  pending: "等待处理",
  accepted: "已退款",
  rejected: "待审核",
  auto_refunded: "超时退款",
  admin_buyer_win: "买家胜",
  admin_seller_win: "卖家胜",
  closed: "已关闭",
};

export default function AdminDisputesPage() {
  const [disputes, setDisputes] = useState<AdminDisputeItem[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [deciding, setDeciding] = useState<number | null>(null);
  const [decision, setDecision] = useState("");
  const [note, setNote] = useState("");

  const fetchDisputes = useCallback(async () => {
    try {
      const res = await apiRequest<Paginated<AdminDisputeItem>>(`/api/admin/disputes?page=${page}&page_size=20`);
      if (res.code === 0) {
        setDisputes(res.data.items || []);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => { fetchDisputes(); }, [fetchDisputes]);

  const submitDecision = async (id: number) => {
    if (!decision) return;
    await apiRequest(`/api/admin/disputes/${id}/decide`, { method: "PUT", body: { decision, note } });
    setDeciding(null);
    setDecision("");
    setNote("");
    fetchDisputes();
  };

  return (
    <div className="max-w-6xl mx-auto px-6 py-16">
      <div className="mb-10">
        <span className="section-label flex items-center gap-2 mb-3"><Gavel size={13} /> Dispute Management</span>
        <h1 className="text-3xl font-bold tracking-tight">争议管理</h1>
      </div>
      {loading ? (
        <div className="glass-card h-96 shimmer" />
      ) : disputes.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <AlertCircle size={32} className="mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500">暂无争议</p>
        </div>
      ) : (
        <div className="space-y-4">
          {disputes.map(item => {
            const d = item.dispute;
            return (
              <div key={d.id} className="glass-card p-6">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <span className="font-medium">争议 #{d.id}</span>
                    <span className="badge !text-[10px] badge-gray">{statusLabels[d.status] || d.status}</span>
                  </div>
                  <span className="text-sm font-bold text-amber-300">{formatCHB(d.refund_amount)}</span>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-4">
                  <div>
                    <div className="text-xs text-gray-500 mb-1">买家原因</div>
                    <p className="text-sm">{d.buyer_reason}</p>
                  </div>
                  {d.seller_reason && (
                    <div>
                      <div className="text-xs text-gray-500 mb-1">卖家回复</div>
                      <p className="text-sm">{d.seller_reason}</p>
                    </div>
                  )}
                </div>
                <div className="text-xs text-gray-600">{formatDate(d.created_at)}</div>
                {d.status === "rejected" && (
                  <div className="mt-4 border-t border-white/5 pt-4">
                    {deciding === d.id ? (
                      <div className="space-y-3">
                        <select className="input" value={decision} onChange={e => setDecision(e.target.value)}>
                          <option value="">选择判定...</option>
                          <option value="buyer_win">买家胜（退款）</option>
                          <option value="seller_win">卖家胜（不退款）</option>
                        </select>
                        <input className="input" placeholder="判定说明" value={note} onChange={e => setNote(e.target.value)} />
                        <div className="flex gap-3">
                          <button className="btn-primary !py-2" onClick={() => submitDecision(d.id)} disabled={!decision}>确认判定</button>
                          <button className="btn-ghost !py-2" onClick={() => setDeciding(null)}>取消</button>
                        </div>
                      </div>
                    ) : (
                      <button className="btn-primary !py-2 !px-4 text-sm" onClick={() => setDeciding(d.id)}>进行判定</button>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
