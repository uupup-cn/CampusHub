'use client';

import { useCallback, useEffect, useState } from 'react';
import { Search, UserX, UserCheck, RefreshCcw, ShieldAlert, Users, Wallet, X } from 'lucide-react';
import { apiRequest, formatCHB, type Paginated } from '@/lib/api';
import { useToast } from '@/components/Toast';

interface User {
  id: number;
  discourse_user_id: number;
  username: string;
  balance: number;
  pending_balance?: number;
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
  const toast = useToast();
  const [adjustUser, setAdjustUser] = useState<User | null>(null);
  const [adjustAmount, setAdjustAmount] = useState('');
  const [adjustReason, setAdjustReason] = useState('');
  const [adjusting, setAdjusting] = useState(false);

  const fetchUsers = useCallback(async () => {
    try {
      const params = new URLSearchParams({ page: '1', page_size: '50' });
      if (keyword) params.set('keyword', keyword);
      const res = await apiRequest<Paginated<User>>('/api/admin/users?' + params);
      if (res.code === 0) setUsers(res.data.items);
    } finally {
      setLoading(false);
    }
  }, [keyword]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const freeze = async (id: number) => {
    const res = await apiRequest('/api/admin/users/' + id + '/freeze', { method: 'PUT' });
    if (res.code === 0) { toast('success', '用户已冻结'); fetchUsers(); } else toast('error', res.message);
  };

  const unfreeze = async (id: number) => {
    const res = await apiRequest('/api/admin/users/' + id + '/unfreeze', { method: 'PUT' });
    if (res.code === 0) { toast('success', '用户已解冻'); fetchUsers(); } else toast('error', res.message);
  };

  const submitAdjust = async (direction: 'add' | 'deduct') => {
    if (!adjustUser) return;
    const amount = Number(adjustAmount);
    if (amount <= 0) { toast('error', '请输入有效金额'); return; }
    if (!adjustReason.trim()) { toast('error', '请填写备注说明'); return; }
    setAdjusting(true);
    const res = await apiRequest('/api/admin/users/' + adjustUser.discourse_user_id + '/adjust', {
      method: 'POST',
      body: { amount, direction, reason: adjustReason.trim() },
    });
    setAdjusting(false);
    if (res.code === 0) {
      toast('success', direction === 'add' ? '积分已增加' : '积分已扣除');
      setAdjustUser(null);
      setAdjustAmount('');
      setAdjustReason('');
      fetchUsers();
    } else {
      toast('error', res.message);
    }
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
        <div className="space-y-3">{[1,2,3].map(function(i) { return <div key={i} className="glass-card h-16 shimmer" /> })}</div>
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
              {users.map(function(u) {
                return (
                  <tr key={u.discourse_user_id}>
                    <td>
                      <div className="flex items-center gap-3">
                        <span className="w-8 h-8 rounded-full bg-gradient-to-br from-violet-500/30 to-cyan-500/20 flex items-center justify-center text-xs font-bold">
                          {u.username ? u.username[0].toUpperCase() : '?'}
                        </span>
                        <div>
                          <div className="font-medium text-gray-200">{u.username || ('用户 ' + u.discourse_user_id)}</div>
                          <div className="text-xs text-gray-600">ID: {u.discourse_user_id}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <button
                        className="inline-flex items-center gap-1.5 font-semibold text-amber-300 hover:text-amber-200 transition-colors cursor-pointer"
                        onClick={function() { setAdjustUser(u); setAdjustAmount(''); setAdjustReason(''); }}
                        title="点击调整积分"
                      >
                        <Wallet size={13} className="text-amber-400" />
                        {formatCHB(u.balance)}
                      </button>
                    </td>
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
                            onClick={function() { unfreeze(u.discourse_user_id); }}
                            title="解冻"
                          >
                            <UserCheck size={14} className="text-emerald-400" />
                          </button>
                        ) : (
                          <button
                            className="w-8 h-8 rounded-full bg-red-500/10 border border-red-500/30 flex items-center justify-center hover:bg-red-500/20 transition-all"
                            onClick={function() { freeze(u.discourse_user_id); }}
                            title="冻结"
                          >
                            <UserX size={14} className="text-red-400" />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {adjustUser && (
        <div className="fixed inset-0 z-[200] flex items-center justify-center bg-black/60 backdrop-blur-sm" onClick={function() { setAdjustUser(null); }}>
          <div className="glass-card !rounded-2xl p-8 max-w-md w-[90%]" onClick={function(e) { e.stopPropagation(); }}>
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2 className="text-lg font-bold">积分调整</h2>
                <p className="text-xs text-gray-500 mt-0.5">{adjustUser.username || ('用户 #' + adjustUser.discourse_user_id)}</p>
              </div>
              <button className="w-8 h-8 rounded-full hover:bg-white/10 flex items-center justify-center" onClick={function() { setAdjustUser(null); }}>
                <X size={16} className="text-gray-400" />
              </button>
            </div>

            <div className="mb-5">
              <div className="flex items-center justify-between p-3 rounded-lg bg-white/5 mb-4">
                <span className="text-xs text-gray-500">当前余额</span>
                <span className="font-bold text-amber-300">{formatCHB(adjustUser.balance)}</span>
              </div>

              <label className="block text-xs text-gray-500 mb-1.5">操作金额 (CHB)</label>
              <input
                className="input mb-4"
                type="number"
                min={1}
                placeholder="输入金额"
                value={adjustAmount}
                onChange={function(e) { setAdjustAmount(e.target.value); }}
              />

              <label className="block text-xs text-gray-500 mb-1.5">备注说明 <span className="text-red-400">*</span></label>
              <textarea
                className="input mb-5 resize-none"
                rows={3}
                placeholder="必须填写调整原因"
                value={adjustReason}
                onChange={function(e) { setAdjustReason(e.target.value); }}
              />

              <div className="flex gap-3">
                <button
                  className="btn-primary flex-1 !py-2.5"
                  disabled={adjusting}
                  onClick={function() { submitAdjust('add'); }}
                >
                  增加积分
                </button>
                <button
                  className="flex-1 py-2.5 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 hover:bg-red-500/20 transition-all font-medium"
                  disabled={adjusting}
                  onClick={function() { submitAdjust('deduct'); }}
                >
                  扣除积分
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
