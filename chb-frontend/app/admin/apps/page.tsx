'use client';

import { useEffect, useState } from 'react';
import { Store, Plus, Trash2, KeyRound, RefreshCcw } from 'lucide-react';
import { apiRequest, formatDate, type Paginated } from '@/lib/api';
import { useToast } from '@/components/Toast';

interface App {
  id: number;
  app_name: string;
  app_description: string | null;
  client_id: string;
  client_secret: string;
  min_trust_level: number;
  fee_rate: number;
  bound_user_id: number | null;
  status: string;
  created_at: string;
}

export default function AdminAppsPage() {
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState({ app_name: '', app_description: '', min_trust_level: 1, fee_rate: 10, bound_user_id: '' });
  const [credentials, setCredentials] = useState<{ client_id: string; client_secret: string } | null>(null);
  const toast = useToast();

  useEffect(() => {
    fetchApps();
  }, []);

  async function fetchApps() {
    setLoading(true);
    try {
      const res = await apiRequest<Paginated<App>>('/api/admin/apps?page=1&page_size=50');
      if (res.code === 0) setApps(res.data.items);
    } finally {
      setLoading(false);
    }
  }

  const createApp = async () => {
    if (!form.app_name) { toast('error', '请填写应用名称'); return; }
    const res = await apiRequest('/api/admin/apps', {
      method: 'POST',
      body: {
        app_name: form.app_name,
        app_description: form.app_description,
        redirect_uris: ['http://localhost:3000/callback'],
        scopes: ['profile:read', 'chb:read', 'chb:spend'],
        min_trust_level: Number(form.min_trust_level),
        fee_rate: Number(form.fee_rate),
        bound_user_id: form.bound_user_id ? Number(form.bound_user_id) : null,
      },
    });
    if (res.code === 0) {
      toast('success', '应用创建成功');
      setCredentials({
        client_id: (res.data as { client_id: string }).client_id,
        client_secret: (res.data as { client_secret: string }).client_secret,
      });
      setShowCreate(false);
      setForm({ app_name: '', app_description: '', min_trust_level: 1, fee_rate: 10, bound_user_id: '' });
      fetchApps();
    } else {
      toast('error', res.message);
    }
  };

  const deleteApp = async (id: number) => {
    if (!confirm('确定删除该应用？')) return;
    const res = await apiRequest(`/api/admin/apps/${id}`, { method: 'DELETE' });
    if (res.code === 0) { toast('success', '应用已删除'); fetchApps(); } else toast('error', res.message);
  };

  return (
    <div className="max-w-6xl mx-auto px-6 py-16">
      <div className="flex items-center justify-between mb-10">
        <div>
          <span className="section-label flex items-center gap-2 mb-3">
            <Store size={13} />
            OAuth2 Applications
          </span>
          <h1 className="text-3xl font-bold tracking-tight">应用管理</h1>
        </div>
        <div className="flex gap-3">
          <button className="btn-ghost !px-4 !py-2" onClick={fetchApps}>
            <RefreshCcw size={14} />
          </button>
          <button className="btn-primary !py-2.5" onClick={() => setShowCreate(!showCreate)}>
            <Plus size={15} />
            新建应用
          </button>
        </div>
      </div>

      {showCreate && (
        <div className="glass-card p-6 mb-8 gradient-border">
          <h2 className="font-semibold mb-6">创建 OAuth2 应用</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div>
              <label className="block text-xs text-gray-500 mb-1.5">应用名称 *</label>
              <input className="input" value={form.app_name} onChange={(e) => setForm({ ...form, app_name: e.target.value })} placeholder="例如：校园二手市场" />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1.5">应用描述</label>
              <input className="input" value={form.app_description} onChange={(e) => setForm({ ...form, app_description: e.target.value })} placeholder="应用用途说明" />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1.5">最低授权等级</label>
              <select className="input" value={form.min_trust_level} onChange={(e) => setForm({ ...form, min_trust_level: Number(e.target.value) })}>
                <option value={0}>TL0</option>
                <option value={1}>TL1</option>
                <option value={2}>TL2</option>
                <option value={3}>TL3</option>
                <option value={4}>TL4</option>
              </select>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1.5">手续费率 (%)</label>
              <input className="input" type="number" min={0} max={100} value={form.fee_rate} onChange={(e) => setForm({ ...form, fee_rate: Number(e.target.value) })} />
            </div>
            <div className="md:col-span-2">
              <label className="block text-xs text-gray-500 mb-1.5">绑定用户 ID（积分收益归属）</label>
              <input className="input" value={form.bound_user_id} onChange={(e) => setForm({ ...form, bound_user_id: e.target.value })} placeholder="留空则积分进入官方池" />
            </div>
          </div>
          <button className="btn-primary" onClick={createApp}>确认创建</button>
        </div>
      )}

      {credentials && (
        <div className="glass-card p-6 mb-8 border border-amber-500/30">
          <div className="flex items-center gap-2 mb-3">
            <KeyRound size={15} className="text-amber-400" />
            <h2 className="font-semibold">应用凭据（仅显示一次，请立即保存）</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <div className="text-xs text-gray-500 mb-1.5">Client ID</div>
              <code className="font-mono text-sm text-violet-300 break-all bg-white/5 rounded-lg px-3 py-2 block">
                {credentials.client_id}
              </code>
            </div>
            <div>
              <div className="text-xs text-gray-500 mb-1.5">Client Secret</div>
              <code className="font-mono text-sm text-amber-300 break-all bg-white/5 rounded-lg px-3 py-2 block">
                {credentials.client_secret}
              </code>
            </div>
          </div>
          <button className="btn-ghost !px-4 !py-2 text-sm" onClick={() => setCredentials(null)}>
            我已保存
          </button>
        </div>
      )}

      {loading ? (
        <div className="space-y-3">{[1,2,3].map(i => <div key={i} className="glass-card h-16 shimmer" />)}</div>
      ) : apps.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Store size={40} className="mx-auto mb-4 text-gray-600" />
          <p className="text-gray-500">暂无应用，点击右上角创建</p>
        </div>
      ) : (
        <div className="space-y-4">
          {apps.map((app) => (
            <div key={app.id} className="glass-card p-5 flex items-center gap-5">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-violet-500/20 to-cyan-500/10 flex items-center justify-center shrink-0">
                <KeyRound size={17} className="text-violet-300" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-3 mb-1">
                  <h3 className="font-semibold">{app.app_name}</h3>
                  <span className={`badge ${app.status === 'active' ? 'badge-green' : 'badge-red'}`}>
                    {app.status === 'active' ? '启用' : '禁用'}
                  </span>
                </div>
                <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500">
                  <span className="font-mono">{app.client_id}</span>
                  <span>最低等级 TL{app.min_trust_level}</span>
                  <span>手续费 {app.fee_rate}%</span>
                  {app.bound_user_id && <span>绑定用户 #{app.bound_user_id}</span>}
                  <span>{formatDate(app.created_at)}</span>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <button
                  className="w-8 h-8 rounded-full bg-red-500/10 border border-red-500/30 flex items-center justify-center hover:bg-red-500/20 transition-all"
                  onClick={() => deleteApp(app.id)}
                  title="删除"
                >
                  <Trash2 size={14} className="text-red-400" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
