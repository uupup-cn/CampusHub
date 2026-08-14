'use client';

import { useCallback, useEffect, useState } from 'react';
import { Wallet } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';

interface Transaction {
  id: number;
  tx_type: string;
  amount: number;
  fee: number;
  net_amount: number;
  from_type: string;
  to_type: string;
  ref_type: string | null;
  description: string | null;
  status: string;
  created_at: string;
}

const typeLabels: Record<string, string> = {
  reward: '奖励',
  spend: '消费',
  recover: '追回',
  release: '释放',
  transfer: '交易',
  admin_adjust: '管理调整',
  refund: '退款',
  refund_deduction: '退款扣减',
};

const refLabels: Record<string, string> = {
  topic: '发帖奖励',
  reply: '回复奖励',
  checkin: '签到奖励',
  like: '被点赞奖励',
};

function getDescription(tx: Transaction): string {
  if (tx.description) return tx.description;
  if (tx.ref_type && refLabels[tx.ref_type]) return refLabels[tx.ref_type];
  return typeLabels[tx.tx_type] || tx.tx_type;
}

function getDisplayAmount(tx: Transaction): number {
  // spend and recover: always negative
  if (tx.tx_type === 'spend' || tx.tx_type === 'recover') {
    return -Math.abs(tx.amount);
  }
  // admin_adjust: amount already has correct sign (positive for add, negative for deduct)
  if (tx.tx_type === 'admin_adjust') {
    return tx.amount;
  }
  // refund_deduction: negative (seller pending deduction)
  if (tx.tx_type === 'refund_deduction') {
    return -Math.abs(tx.amount);
  }
  // reward, release, transfer, refund: positive
  return Math.abs(tx.amount);
}

export default function DashboardPointsPage() {
  const [txs, setTxs] = useState<Transaction[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [filter, setFilter] = useState<string>('');

  const fetchTxs = useCallback(async () => {
    try {
      const params = new URLSearchParams({ page: String(page), page_size: '20' });
      if (filter) params.set('type', filter);
      const res = await apiRequest<Paginated<Transaction>>('/api/chb/me/transactions?' + params);
      if (res.code === 0) {
        setTxs(res.data.items);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  }, [page, filter]);

  useEffect(() => { fetchTxs(); }, [fetchTxs]);

  const totalPages = Math.max(1, Math.ceil(total / 20));

  const filters = [
    { v: '', l: '全部' },
    { v: 'reward', l: '奖励' },
    { v: 'spend', l: '消费' },
    { v: 'transfer', l: '交易' },
    { v: 'admin_adjust', l: '管理调整' },
    { v: 'recover', l: '追回' },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">积分明细</h1>

      <div className="flex gap-2 mb-6">
        {filters.map(f => (
          <button
            key={f.v}
            onClick={() => { setFilter(f.v); setPage(1); }}
            className={'px-3 py-1.5 rounded-full text-xs transition-all ' + (filter === f.v ? 'bg-violet-500/20 text-violet-300 border border-violet-500/30' : 'text-gray-400 border border-white/5')}
          >
            {f.l}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="space-y-2">
          {[1, 2, 3, 4, 5].map(function(i) { return <div key={i} className="glass-card h-14 shimmer" /> })}
        </div>
      ) : txs.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <Wallet size={32} className="mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500">暂无积分记录</p>
        </div>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="data-table">
            <thead>
              <tr>
                <th>类型</th>
                <th>描述</th>
                <th>金额</th>
                <th>时间</th>
              </tr>
            </thead>
            <tbody>
              {txs.map(function(tx) {
                const displayAmount = getDisplayAmount(tx);
                return (
                  <tr key={tx.id}>
                    <td><span className="badge !text-[10px] badge-gray">{typeLabels[tx.tx_type] || tx.tx_type}</span></td>
                    <td className="truncate max-w-xs">{getDescription(tx)}</td>
                    <td className={displayAmount > 0 ? 'text-emerald-300 font-medium' : 'text-red-300 font-medium'}>
                      {displayAmount > 0 ? '+' : ''}{formatCHB(displayAmount)}
                    </td>
                    <td className="text-gray-600 text-xs">{formatDate(tx.created_at)}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-3 mt-8">
          <button className="btn-ghost !px-4 !py-2 text-sm" disabled={page <= 1} onClick={function() { setPage(page - 1); }}>上一页</button>
          <span className="text-sm text-gray-500">{page} / {totalPages}</span>
          <button className="btn-ghost !px-4 !py-2 text-sm" disabled={page >= totalPages} onClick={function() { setPage(page + 1); }}>下一页</button>
        </div>
      )}
    </div>
  );
}
