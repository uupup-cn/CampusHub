'use client';

import { Suspense, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { Loader2, CheckCircle2, AlertCircle } from 'lucide-react';
import { useState } from 'react';

function CallbackContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');

  useEffect(() => {
    // Discourse 登录成功后会带 return_url 跳回来
    // 检查 /api/auth/me 确认登录状态
    fetch('/api/auth/me').then(r => r.json()).then(data => {
      if (data.code === 0 && data.data?.logged_in) {
        localStorage.setItem('chb_user_id', String(data.data.user_id));
        setStatus('success');
        setMessage('欢迎回来，' + data.data.username);
        setTimeout(() => router.push('/'), 1500);
      } else {
        setStatus('error');
        setMessage('未检测到登录状态，请重新登录');
        setTimeout(() => router.push('/'), 3000);
      }
    }).catch(() => {
      setStatus('error');
      setMessage('网络错误');
      setTimeout(() => router.push('/'), 3000);
    });
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center px-6 py-16">
      <div className="w-full max-w-md">
        <div className="glass-card p-10 text-center">
          {status === 'loading' && (
            <div className="flex flex-col items-center gap-4">
              <Loader2 size={32} className="text-violet-400 animate-spin" />
              <p className="text-gray-400 text-sm">正在同步登录状态...</p>
            </div>
          )}
          {status === 'success' && (
            <div className="flex flex-col items-center gap-4">
              <CheckCircle2 size={40} className="text-emerald-400" />
              <h1 className="text-xl font-bold">登录成功</h1>
              <p className="text-gray-500 text-sm text-center">{message}</p>
            </div>
          )}
          {status === 'error' && (
            <div className="flex flex-col items-center gap-4">
              <AlertCircle size={40} className="text-rose-400" />
              <h1 className="text-xl font-bold">登录失败</h1>
              <p className="text-gray-500 text-sm text-center">{message}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default function CallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 size={32} className="text-violet-400 animate-spin" />
      </div>
    }>
      <CallbackContent />
    </Suspense>
  );
}
