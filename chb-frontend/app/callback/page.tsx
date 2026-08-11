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
      setMessage('Missing authorization code');
      return;
    }

    // The authorization code is available in the URL.
    // The calling application should exchange it for a token via POST /oauth/token.
    // This page simply displays the result and provides the code to the parent window.
    setStatus('success');
    setMessage('Authorization successful. Code: ' + code.substring(0, 8) + '...');

    // If this is a popup window, send the code back to the opener
    if (window.opener) {
      window.opener.postMessage(
        { type: 'oauth_callback', code, state },
        '*'
      );
      setTimeout(() => window.close(), 2000);
    }
  }, [searchParams, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#0a0a0f]">
      <div className="max-w-md w-full mx-4 p-8 rounded-2xl bg-[#111118] border border-white/10">
        {status === 'loading' && (
          <div className="flex flex-col items-center gap-4">
            <Loader2 className="w-12 h-12 text-indigo-400 animate-spin" />
            <p className="text-white/60 text-sm">Processing authorization...</p>
          </div>
        )}
        {status === 'success' && (
          <div className="flex flex-col items-center gap-4">
            <CheckCircle2 className="w-12 h-12 text-emerald-400" />
            <h1 className="text-xl font-bold text-white">Authorization Successful</h1>
            <p className="text-white/50 text-sm text-center">{message}</p>
            <button
              onClick={() => router.push('/')}
              className="mt-4 px-6 py-2 rounded-lg bg-indigo-500 hover:bg-indigo-400 text-white text-sm font-medium transition-colors"
            >
              Return Home
            </button>
          </div>
        )}
        {status === 'error' && (
          <div className="flex flex-col items-center gap-4">
            <AlertCircle className="w-12 h-12 text-rose-400" />
            <h1 className="text-xl font-bold text-white">Authorization Failed</h1>
            <p className="text-white/50 text-sm text-center">{message}</p>
            <button
              onClick={() => router.push('/')}
              className="mt-4 px-6 py-2 rounded-lg bg-white/10 hover:bg-white/20 text-white text-sm font-medium transition-colors"
            >
              Return Home
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default function CallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-[#0a0a0f]">
        <Loader2 className="w-8 h-8 text-indigo-400 animate-spin" />
      </div>
    }>
      <CallbackContent />
    </Suspense>
  );
}
