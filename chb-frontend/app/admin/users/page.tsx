'use client';

import { useCallback, useEffect, useState } from 'react';
import { Search, UserX, UserCheck, RefreshCcw, ShieldAlert, Users } from 'lucide-react';
import { apiRequest, formatCHB, type Paginated } from '@/lib/api';
import { useToast } from '@/components/Toast';

interface User {
  id: number;
  discourse_user_id: number;
  username: string;
  balance: number;
  trust_level: number;
  total_earned: number;
  total_spent: number;
  status: string;
  created_at: string;
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [keyword, setKeyword] = useState('');
  const [recoverAmount, setRecoverAmount] = useState<Record<number, string>>({});
  const toast = useToast();

  const fetchUsers = useCallback(async () => {
    try {
      const params = new URLSearchParams({ page: '1', page_size: '50' });
      if (keyword) params.set('keyword', keyword);
      const res = await apiRequest<Paginated<User>>(`/api/admin/users?${params}`);
      if (res.code === 0) setUsers(res.data.items);
    } finally {
      setLoading(false);
    }
  }, [keyword]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const freeze = async (id: number) => {
    const res = await apiRequest(`/api/admin/users/${id}/freeze`, { method: 'PUT' });
    if (res.code === 0) { toast('success', '用户已冻结'); fetchUsers(); } else toast('error', res.message);
  };

  const unfreeze = async (id: number) => {
    const res = await apiRequest(`/api/admin/users/${id}/unfreeze`, { method: 'PUT' });
    if (res.code === 0) { toast('success', '用户已解冻'); fetchUsers(); } else toast('error', res.message);
  };

  const recover = async (userId: number) => {
    const amount = Number(recoverAmount[userId] || 0);
    if (amount <= 0) { toast('error', '请输入有效的追回金额'); return; }
    const res = await apiRequest(`/api/admin/users/${userId}/recover`, {
      method: 'POST',
      body: { amount, reason: '管理员追回' },
    });
    if (res.code === 0) { toast('success', '积分已追回'); fetchUsers(); } else toast('error', res.message);
  };

  return (
    <div className="max-w-6xl mx-auto px-6 py-16">
      <div className="flex items-center justify-between mb-10">
        <div>
          <span className="section-label flex items-center gap-2 mb-3">
            <Users size={13} />
            User Management
          </span>
          <h1 className="text-3xl font-bold tracking-tight">用户管理</h1>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative">
            <Search size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-500" />
            <input
              className="input !py-2 !pl-9 !pr-4 w-56 text-sm"
              placeholder="搜索用户名..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
            />
          </div>
          <button className="btn-ghost !px-3 !py-2" onClick={fetchUsers}>
            <RefreshCcw size={14} />
          </button>
        </div>
      </div>

      {loading ? (
        <div className="space-y-3">{[1,2,3].map(i => <div key={i} className="glass-card h-16 shimmer" />)}</div>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="data-table">
            <thead>
              <tr>
                <th>用户</th>
                <th>余额</th>
                <th>等级</th>
                <th>累计赚取</th>
                <th>累计消费</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.discourse_user_id}>
                  <td>
                    <div className="flex items-center gap-3">
                      <span className="w-8 h-8 rounded-full bg-gradient-to-br from-violet-500/30 to-cyan-500/20 flex items-center justify-center text-xs font-bold">
                        {u.username?.[0]?.toUpperCase() || '?'}
                      </span>
                      <div>
                        <div className="font-medium text-gray-200">{u.username || `用户 ${u.discourse_user_id}`}</div>
                        <div className="text-xs text-gray-600">ID: {u.discourse_user_id}</div>
                      </div>
                    </div>
                  </td>
                  <td className="font-semibold text-gray-200">{formatCHB(u.balance)}</td>
                  <td><span className="badge badge-purple">TL{u.trust_level}</span></td>
                  <td className="text-emerald-400">{formatCHB(u.total_earned)}</td>
                  <td className="text-amber-400">{formatCHB(u.total_spent)}</td>
                  <td>
                    {u.status === 'frozen' ? (
                      <span className="badge badge-red"><ShieldAlert size={12} />已冻结</span>
                    ) : (
                      <span className="badge badge-green">正常</span>
                    )}
                  </td>
                  <td>
                    <div className="flex items-center gap-2">
                      {u.status === 'frozen' ? (
                        <button
                          className="w-8 h-8 rounded-full bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center hover:bg-emerald-500/20 transition-all"
                          onClick={() => unfreeze(u.discourse_user_id)}
                          title="解冻"
                        >
                          <UserCheck size={14} className="text-emerald-400" />
                        </button>
                      ) : (
                        <button
                          className="w-8 h-8 rounded-full bg-red-500/10 border border-red-500/30 flex items-center justify-center hover:bg-red-500/20 transition-all"
                          onClick={() => freeze(u.discourse_user_id)}
                          title="冻结"
                        >
                          <UserX size={14} className="text-red-400" />
                        </button>
                      )}
                      <input
                        className="!py-1 !px-2 !text-xs w-20 input"
                        placeholder="追回"
                        value={recoverAmount[u.discourse_user_id] || ''}
                        onChange={(e) => setRecoverAmount({ ...recoverAmount, [u.discourse_user_id]: e.target.value })}
                      />
                      <button
                        className="px-2.5 py-1 rounded-lg bg-amber-500/10 border border-amber-500/30 text-[11px] text-amber-400 hover:bg-amber-500/20 transition-all"
                        onClick={() => recover(u.discourse_user_id)}
                      >
                        追回
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
