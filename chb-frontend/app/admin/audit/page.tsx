'use client';

import { useEffect, useState } from 'react';
import { FileText, Search, RefreshCcw } from 'lucide-react';
import { apiRequest, formatDate, type Paginated } from '@/lib/api';

interface AuditLog {
  id: number;
  operator_id: number | null;
  action: string;
  target_type: string | null;
  target_id: number | null;
  detail: string | null;
  ip_address: string | null;
  created_at: string;
}

export default function AdminAuditPage() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [action, setAction] = useState('');

  useEffect(() => {
    fetchLogs();
  }, [action]);

  async function fetchLogs() {
    setLoading(true);
    try {
      const params = new URLSearchParams({ page: '1', page_size: '50' });
      if (action) params.set('action', action);
      const res = await apiRequest<Paginated<AuditLog>>(`/api/admin/audit-logs?${params}`);
      if (res.code === 0) setLogs(res.data.items);
    } finally {
      setLoading(false);
    }
  }

  const actionLabel = (a: string) => {
    const map: Record<string, string> = {
      config_update: '配置更新', app_create: '创建应用', app_update: '更新应用', app_delete: '删除应用',
      merchant_approve: '审核入驻', item_review: '审核商品', user_freeze: '冻结用户', user_unfreeze: '解冻用户',
      points_recover: '积分追回', release_manual: '手动释放',
    };
    return map[a] || a;
  };

  return (
    <div className="max-w-6xl mx-auto px-6 py-16">
      <div className="flex items-center justify-between mb-10">
        <div>
          <span className="section-label flex items-center gap-2 mb-3">
            <FileText size={13} />
            Audit Logs
          </span>
          <h1 className="text-3xl font-bold tracking-tight">审计日志</h1>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative">
            <Search size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-500" />
            <input
              className="input !py-2 !pl-9 !pr-4 w-64 text-sm"
              placeholder="筛选操作类型..."
              value={action}
              onChange={(e) => setAction(e.target.value)}
            />
          </div>
          <button className="btn-ghost !px-3 !py-2" onClick={fetchLogs}>
            <RefreshCcw size={14} />
          </button>
        </div>
      </div>

      {loading ? (
        <div className="space-y-3">{[1,2,3].map(i => <div key={i} className="glass-card h-14 shimmer" />)}</div>
      ) : logs.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <FileText size={40} className="mx-auto mb-4 text-gray-600" />
          <p className="text-gray-500">暂无审计记录</p>
        </div>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="data-table">
            <thead>
              <tr>
                <th>操作</th>
                <th>操作人</th>
                <th>目标</th>
                <th>IP</th>
                <th>时间</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((log) => (
                <tr key={log.id}>
                  <td><span className="badge badge-purple">{actionLabel(log.action)}</span></td>
                  <td>{log.operator_id ? `#${log.operator_id}` : '-'}</td>
                  <td className="text-xs">
                    {log.target_type ? `${log.target_type} #${log.target_id}` : '-'}
                  </td>
                  <td className="font-mono text-xs">{log.ip_address || '-'}</td>
                  <td className="text-xs">{formatDate(log.created_at)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}