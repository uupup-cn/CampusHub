'use client';

import { useEffect, useState } from 'react';
import { Settings, Save, Trophy, Award, Coins } from 'lucide-react';
import { apiRequest } from '@/lib/api';
import { useToast } from '@/components/Toast';

interface RewardRule {
  action: string;
  display_name: string;
  amount: number;
  cooldown_seconds: number;
  daily_cap_per_user: number;
  is_active: boolean;
}

interface TrustLevel {
  trust_level: number;
  daily_cap: number;
  reward_multiplier: number;
}

interface SystemSettings {
  marketplace_fee_rate: number;
  auto_release_enabled: boolean;
  auto_release_threshold: number;
  auto_release_ratio: number;
  auto_release_monthly_cap: number;
}

export default function AdminSettingsPage() {
  const [rules, setRules] = useState<RewardRule[]>([]);
  const [caps, setCaps] = useState<TrustLevel[]>([]);
  const [settings, setSettings] = useState<SystemSettings>({
    marketplace_fee_rate: 10,
    auto_release_enabled: true,
    auto_release_threshold: 80,
    auto_release_ratio: 50,
    auto_release_monthly_cap: 10000000,
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const toast = useToast();

  useEffect(() => {
    Promise.all([
      apiRequest<RewardRule[]>('/api/chb/reward/rules'),
      apiRequest<TrustLevel[]>('/api/admin/trust-levels'),
      apiRequest<SystemSettings>('/api/admin/settings'),
    ]).then(([r, c, s]) => {
      if (r.code === 0) setRules(r.data);
      if (c.code === 0) setCaps(c.data);
      if (s.code === 0) setSettings(s.data);
    }).finally(() => setLoading(false));
  }, []);

  const updateRule = async (action: string, field: string, value: number | boolean) => {
    const rule = rules.find((r) => r.action === action);
    if (!rule) return;
    const updated = { ...rule, [field]: value };
    const res = await apiRequest('/api/admin/reward/rules', {
      method: 'PUT',
      body: { action, amount: updated.amount, cooldown_seconds: updated.cooldown_seconds, daily_cap_per_user: updated.daily_cap_per_user, is_active: updated.is_active },
    });
    if (res.code === 0) {
      setRules(rules.map((r) => (r.action === action ? updated : r)));
      toast('success', `${rule.display_name}规则已更新`);
    } else {
      toast('error', res.message);
    }
  };

  const updateCap = async (level: number, field: string, value: number) => {
    const cap = caps.find((c) => c.trust_level === level);
    if (!cap) return;
    const updated = { ...cap, [field]: value };
    const res = await apiRequest('/api/admin/trust-levels', {
      method: 'PUT',
      body: { trust_level: level, daily_cap: updated.daily_cap, reward_multiplier: updated.reward_multiplier },
    });
    if (res.code === 0) {
      setCaps(caps.map((c) => (c.trust_level === level ? updated : c)));
      toast('success', `TL${level} 已更新`);
    } else {
      toast('error', res.message);
    }
  };

  const saveSettings = async () => {
    setSaving(true);
    try {
      const res = await apiRequest('/api/admin/settings', {
        method: 'PUT',
        body: settings,
      });
      if (res.code === 0) toast('success', '系统配置已保存');
      else toast('error', res.message);
    } finally {
      setSaving(false);
    }
  };

  const actionLabels: Record<string, string> = { post: '发帖', reply: '回复', checkin: '签到', liked: '被点赞' };

  if (loading) {
    return <div className="max-w-6xl mx-auto px-6 py-16"><div className="glass-card h-96 shimmer" /></div>;
  }

  return (
    <div className="max-w-6xl mx-auto px-6 py-16">
      <div className="mb-10">
        <span className="section-label flex items-center gap-2 mb-3">
          <Settings size={13} />
          System Configuration
        </span>
        <h1 className="text-3xl font-bold tracking-tight">系统配置</h1>
      </div>

      <div className="space-y-8">
        {/* Reward Rules */}
        <section className="glass-card p-6">
          <h2 className="flex items-center gap-2 font-semibold mb-6">
            <Trophy size={17} className="text-violet-400" />
            奖励规则
          </h2>
          <div className="space-y-4">
            {rules.map((rule) => (
              <div key={rule.action} className="flex items-center gap-4 p-4 rounded-xl bg-white/2 border border-white/5">
                <div className="w-16 shrink-0">
                  <span className={`badge ${rule.is_active ? 'badge-green' : 'badge-gray'}`}>
                    {actionLabels[rule.action] || rule.action}
                  </span>
                </div>
                <div className="flex items-center gap-2 flex-1">
                  <label className="text-xs text-gray-500 w-16">每次奖励</label>
                  <input
                    type="number"
                    className="input !py-1.5 w-28 text-sm"
                    value={rule.amount}
                    onChange={(e) => updateRule(rule.action, 'amount', Number(e.target.value))}
                  />
                </div>
                <div className="flex items-center gap-2">
                  <label className="text-xs text-gray-500 w-20">冷却(秒)</label>
                  <input
                    type="number"
                    className="input !py-1.5 w-28 text-sm"
                    value={rule.cooldown_seconds}
                    onChange={(e) => updateRule(rule.action, 'cooldown_seconds', Number(e.target.value))}
                  />
                </div>
                <div className="flex items-center gap-2">
                  <label className="text-xs text-gray-500 w-24">每日上限</label>
                  <input
                    type="number"
                    className="input !py-1.5 w-28 text-sm"
                    value={rule.daily_cap_per_user}
                    onChange={(e) => updateRule(rule.action, 'daily_cap_per_user', Number(e.target.value))}
                  />
                </div>
                <button
                  className={`px-3 py-1.5 rounded-full text-xs transition-all ${
                    rule.is_active
                      ? 'bg-emerald-500/15 text-emerald-400 border border-emerald-500/30'
                      : 'bg-white/5 text-gray-500 border border-white/10'
                  }`}
                  onClick={() => updateRule(rule.action, 'is_active', !rule.is_active)}
                >
                  {rule.is_active ? '启用中' : '已停用'}
                </button>
              </div>
            ))}
          </div>
        </section>

        {/* Trust Level Caps */}
        <section className="glass-card p-6">
          <h2 className="flex items-center gap-2 font-semibold mb-6">
            <Award size={17} className="text-cyan-400" />
            等级日上限
          </h2>
          <div className="space-y-4">
            {caps.map((cap) => (
              <div key={cap.trust_level} className="flex items-center gap-4 p-4 rounded-xl bg-white/2 border border-white/5">
                <span className="badge badge-purple w-16 justify-center">TL{cap.trust_level}</span>
                <div className="flex items-center gap-2 flex-1">
                  <label className="text-xs text-gray-500 w-24">每日上限</label>
                  <input
                    type="number"
                    className="input !py-1.5 w-36 text-sm"
                    value={cap.daily_cap}
                    onChange={(e) => updateCap(cap.trust_level, 'daily_cap', Number(e.target.value))}
                  />
                </div>
                <div className="flex items-center gap-2">
                  <label className="text-xs text-gray-500 w-24">奖励倍率</label>
                  <input
                    type="number"
                    step="0.01"
                    className="input !py-1.5 w-28 text-sm"
                    value={cap.reward_multiplier}
                    onChange={(e) => updateCap(cap.trust_level, 'reward_multiplier', Number(e.target.value))}
                  />
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* System Settings */}
        <section className="glass-card p-6">
          <h2 className="flex items-center gap-2 font-semibold mb-6">
            <Coins size={17} className="text-amber-400" />
            经济参数
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
            <div>
              <label className="block text-xs text-gray-500 mb-2">集市手续费率 (%)</label>
              <input
                type="number"
                className="input"
                value={settings.marketplace_fee_rate}
                onChange={(e) => setSettings({ ...settings, marketplace_fee_rate: Number(e.target.value) })}
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-2">自动释放阈值 (%)</label>
              <input
                type="number"
                className="input"
                value={settings.auto_release_threshold}
                onChange={(e) => setSettings({ ...settings, auto_release_threshold: Number(e.target.value) })}
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-2">自动释放比例 (%)</label>
              <input
                type="number"
                className="input"
                value={settings.auto_release_ratio}
                onChange={(e) => setSettings({ ...settings, auto_release_ratio: Number(e.target.value) })}
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-2">月度释放上限 (CHB)</label>
              <input
                type="number"
                className="input"
                value={settings.auto_release_monthly_cap}
                onChange={(e) => setSettings({ ...settings, auto_release_monthly_cap: Number(e.target.value) })}
              />
            </div>
            <div className="md:col-span-2">
              <label className="flex items-center gap-3 cursor-pointer">
                <button
                  className={`w-11 h-6 rounded-full relative transition-colors ${settings.auto_release_enabled ? 'bg-violet-500' : 'bg-white/10'}`}
                  onClick={() => setSettings({ ...settings, auto_release_enabled: !settings.auto_release_enabled })}
                >
                  <span
                    className={`absolute top-0.5 w-5 h-5 rounded-full bg-white transition-all ${settings.auto_release_enabled ? 'left-[22px]' : 'left-0.5'}`}
                  />
                </button>
                <span className="text-sm text-gray-300">自动释放阀</span>
                <span className="text-xs text-gray-600">{settings.auto_release_enabled ? '公共池水位低于阈值时自动回补' : '仅手动释放'}</span>
              </label>
            </div>
          </div>
          <button className="btn-primary" onClick={saveSettings} disabled={saving}>
            <Save size={15} />
            {saving ? '保存中...' : '保存配置'}
          </button>
        </section>
      </div>
    </div>
  );
}