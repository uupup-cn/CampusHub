'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Store, Loader2 } from 'lucide-react';
import { apiRequest } from '@/lib/api';

export default function DashboardApplyPage() {
  const router = useRouter();
  const [shopName, setShopName] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!shopName.trim()) {
      setError('请输入店铺名称');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const res = await apiRequest('/api/marketplace/apply', {
        method: 'POST',
        body: { shop_name: shopName, description },
      });
      if (res.code === 0) {
        router.push('/dashboard/marketplace');
      } else {
        setError(res.message);
      }
    } catch {
      setError('提交失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-8">申请入驻集市</h1>

      <form onSubmit={handleSubmit} className="glass-card p-8 max-w-lg">
        <div className="mb-6">
          <label className="block text-sm font-medium mb-2">店铺名称</label>
          <input
            className="input"
            placeholder="给你的店铺起个名字"
            value={shopName}
            onChange={(e) => setShopName(e.target.value)}
          />
        </div>
        <div className="mb-6">
          <label className="block text-sm font-medium mb-2">店铺描述</label>
          <textarea
            className="input"
            placeholder="简述你的店铺定位和商品类型"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={4}
          />
        </div>
        {error && <p className="text-sm text-red-400 mb-4">{error}</p>}
        <button type="submit" className="btn-primary w-full" disabled={loading}>
          {loading ? <Loader2 size={16} className="animate-spin" /> : <Store size={16} />}
          提交申请
        </button>
      </form>
    </div>
  );
}
