'use client';

import { useCallback, useEffect, useState } from 'react';
import { Wallet } from 'lucide-react';
import { apiRequest, formatCHB, formatDate, type Paginated } from '@/lib/api';

interface Transaction {
  id: number;
  type: string;
  amount: number;
  balance_after: number;
  description: string;
  ref_type: string;
  created_at: string;
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
      const res = await apiRequest<Paginated<Transaction>>(`/api/chb/me/transactions?${params}`);
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

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">积分明细</h1>

      <div className="flex gap-2 mb-6">
        {[{ v: '', l: '全部' }, { v: 'reward', l: '奖励' }, { v: 'spend', l: '消费' }, { v: 'recover', l: '追回' }].map(f => (
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
          {[1, 2, 3, 4, 5].map(i => <div key={i} className="glass-card h-14 shimmer" />)}
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
                <th>余额</th>
                <th>时间</th>
              </tr>
            </thead>
            <tbody>
              {txs.map(tx => (
                <tr key={tx.id}>
                  <td><span className="badge !text-[10px] badge-gray">{tx.type}</span></td>
                  <td className="truncate max-w-xs">{tx.description}</td>
                  <td className={tx.amount > 0 ? 'text-emerald-300 font-medium' : 'text-red-300 font-medium'}>
                    {tx.amount > 0 ? '+' : ''}{formatCHB(tx.amount)}
                  </td>
                  <td className="text-gray-500">{formatCHB(tx.balance_after)}</td>
                  <td className="text-gray-600 text-xs">{formatDate(tx.created_at)}</td>
                </tr>
              ))}
            </tbody>
          </table>
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
