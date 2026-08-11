'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import { useState, useEffect } from 'react';
import { Loader2, ShieldCheck, UserCircle2, Coins, ShoppingBag, CheckCircle2 } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface AppInfo {
  app_name: string;
  app_description: string;
  scopes: string[];
  min_trust_level: number;
  redirect_uris: string[];
}

const scopeLabels: Record<string, string> = {
  'profile:read': '读取你的个人资料',
  'chb:read': '查询你的 CHB 积分余额',
  'chb:spend': '扣减你的 CHB 积分',
  'forum:read': '读取论坛公开内容',
  'marketplace:read': '查看集市商品',
  'marketplace:trade': '在集市进行交易',
};

const scopeIcons: Record<string, React.ReactNode> = {
  'profile:read': <UserCircle2 size={18} className="text-cyan-400" />,
  'chb:read': <Coins size={18} className="text-amber-400" />,
  'chb:spend': <ShoppingBag size={18} className="text-violet-400" />,
  'forum:read': <ShieldCheck size={18} className="text-emerald-400" />,
};

function AuthorizeContent() {
  const searchParams = useSearchParams();
  const [appInfo, setAppInfo] = useState<AppInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const clientId = searchParams.get('client_id') || '';
  const redirectUri = searchParams.get('redirect_uri') || '';
  const scope = searchParams.get('scope') || '';
  const state = searchParams.get('state') || '';

  useEffect(() => {
    async function fetchAppInfo() {
      try {
        const res = await fetch(`/api/oauth/app-info?client_id=${clientId}`);
        const data = await res.json();
        if (data.code === 0) {
          setAppInfo(data.data);
        } else {
          setError('无法获取应用信息');
        }
      } catch {
        setError('网络错误');
      } finally {
        setLoading(false);
      }
    }
    if (clientId) fetchAppInfo();
  }, [clientId]);

  const [submitting, setSubmitting] = useState(false);

  const handleAuthorize = async () => {
    if (submitting) return;
    setSubmitting(true);
    try {
      const res = await apiRequest<{ redirect_uri: string }>('/api/oauth/authorize/confirm', {
        method: 'POST',
        body: {
          client_id: clientId,
          redirect_uri: redirectUri,
          response_type: 'code',
          scope,
          state,
        },
      });
      if (res.code === 0 && res.data?.redirect_uri) {
        window.location.assign(res.data.redirect_uri);
      } else {
        setError(res.message || '授权失败');
        setSubmitting(false);
      }
    } catch {
      setError('网络错误');
      setSubmitting(false);
    }
  };

  const handleCancel = () => {
    window.location.assign(redirectUri + '?error=access_denied&state=' + state);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 size={32} className="text-violet-400 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="glass-card p-10 text-center">
          <p className="text-red-400 mb-4">{error}</p>
          <button className="btn-ghost" onClick={handleCancel}>返回</button>
        </div>
      </div>
    );
  }

  const scopesList = scope.split(' ').filter(Boolean);

  return (
    <div className="min-h-screen flex items-center justify-center px-6 py-16">
      <div className="w-full max-w-md">
        <div className="text-center mb-10">
          <span className="section-label">OAuth2 Authorization</span>
          <h1 className="display-title text-4xl mt-3 mb-4">授权请求</h1>
          <p className="text-gray-400 text-sm">
            应用 <span className="text-white font-semibold">{appInfo?.app_name || clientId}</span> 请求以下权限
          </p>
        </div>

        <div className="glass-card p-8 mb-6">
          <div className="space-y-4">
            {scopesList.map((s) => (
              <div key={s} className="flex items-start gap-3 p-3 rounded-xl bg-white/3 border border-white/5">
                <div className="shrink-0 mt-0.5">{scopeIcons[s] || <ShieldCheck size={18} className="text-gray-500" />}</div>
                <div>
                  <div className="text-sm font-medium text-gray-200">{scopeLabels[s] || s}</div>
                  <div className="text-xs text-gray-500 mt-0.5 font-mono">{s}</div>
                </div>
              </div>
            ))}
          </div>

          {appInfo?.app_description && (
            <p className="text-sm text-gray-500 mt-6 leading-relaxed">{appInfo.app_description}</p>
          )}
        </div>

        <div className="flex gap-4">
          <button className="btn-ghost flex-1" onClick={handleCancel}>取消</button>
          <button className="btn-primary flex-1" onClick={handleAuthorize} disabled={submitting}>
            <CheckCircle2 size={16} />
            {submitting ? '处理中...' : '同意授权'}
          </button>
        </div>
      </div>
    </div>
  );
}

export default function AuthorizePage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 size={32} className="text-violet-400 animate-spin" />
      </div>
    }>
      <AuthorizeContent />
    </Suspense>
  );
}
