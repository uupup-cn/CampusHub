'use client';

import { Suspense, useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { CheckCircle2, AlertCircle, Loader2 } from 'lucide-react';

function CallbackContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state');
    const error = searchParams.get('error');

    if (error) {
      setStatus('error');
      setMessage(decodeURIComponent(error));
      return;
    }

    if (!code) {
      setStatus('error');
      setMessage('缺少授权码');
      return;
    }

    setStatus('success');
    setMessage('授权成功，授权码：' + code.substring(0, 8) + '...');

    if (window.opener) {
      window.opener.postMessage(
        { type: 'oauth_callback', code, state },
        '*'
      );
      setTimeout(() => window.close(), 2000);
    }
  }, [searchParams, router]);

  return (
    <div className="min-h-screen flex items-center justify-center px-6 py-16">
      <div className="w-full max-w-md">
        <div className="glass-card p-10 text-center">
          {status === 'loading' && (
            <div className="flex flex-col items-center gap-4">
              <Loader2 size={32} className="text-violet-400 animate-spin" />
              <p className="text-gray-400 text-sm">正在处理授权...</p>
            </div>
          )}
          {status === 'success' && (
            <div className="flex flex-col items-center gap-4">
              <CheckCircle2 size={40} className="text-emerald-400" />
              <h1 className="text-xl font-bold">授权成功</h1>
              <p className="text-gray-500 text-sm text-center">{message}</p>
              <button
                onClick={() => router.push('/')}
                className="btn-primary mt-4"
              >
                返回首页
              </button>
            </div>
          )}
          {status === 'error' && (
            <div className="flex flex-col items-center gap-4">
              <AlertCircle size={40} className="text-rose-400" />
              <h1 className="text-xl font-bold">授权失败</h1>
              <p className="text-gray-500 text-sm text-center">{message}</p>
              <button
                onClick={() => router.push('/')}
                className="btn-ghost mt-4"
              >
                返回首页
              </button>
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
