'use client';

import { useEffect, useState } from 'react';
import { Plug } from 'lucide-react';
import { apiRequest, formatDate } from '@/lib/api';

interface App {
  id: number;
  app_name: string;
  client_id: string;
  scopes: string;
  created_at: string;
}

export default function DashboardAppsPage() {
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiRequest<{ items: App[] }>('/api/oauth/my-apps')
      .then(res => { if (res.code === 0) setApps(res.data.items || []); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">授权应用</h1>

      {loading ? (
        <div className="space-y-3">
          {[1, 2].map(i => <div key={i} className="glass-card h-20 shimmer" />)}
        </div>
      ) : apps.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <Plug size={32} className="mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500">暂无已授权的应用</p>
        </div>
      ) : (
        <div className="space-y-3">
          {apps.map(app => (
            <div key={app.id} className="glass-card p-5 flex items-center justify-between">
              <div>
                <div className="font-medium mb-1">{app.app_name}</div>
                <div className="text-xs text-gray-500">Client ID: {app.client_id}</div>
                <div className="text-xs text-gray-600 mt-1">授权时间: {formatDate(app.created_at)}</div>
              </div>
              <div className="flex gap-2">
                {(app.scopes || '').split(',').map(s => (
                  <span key={s} className="badge badge-cyan !text-[10px]">{s}</span>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
