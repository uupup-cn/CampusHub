'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Store, Send, CheckCircle2, AlertCircle } from 'lucide-react';
import { apiRequest } from '@/lib/api';
import { useToast } from '@/components/Toast';

export default function ApplyMerchantPage() {
  const router = useRouter();
  const toast = useToast();
  const [shopName, setShopName] = useState('');
  const [description, setDescription] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async () => {
    if (!shopName.trim()) {
      toast('error', '请填写店铺名称');
      return;
    }
    setSubmitting(true);
    try {
      const res = await apiRequest('/api/marketplace/apply', {
        method: 'POST',
        body: { shop_name: shopName, description },
      });
      if (res.code === 0) {
        setSubmitted(true);
      } else {
        toast('error', res.message);
      }
    } catch {
      toast('error', '网络错误，请稍后重试');
    } finally {
      setSubmitting(false);
    }
  };

  if (submitted) {
    return (
      <div className="max-w-lg mx-auto px-6 py-24 text-center">
        <div className="w-16 h-16 rounded-full bg-emerald-500/15 border border-emerald-500/30 flex items-center justify-center mx-auto mb-6">
          <CheckCircle2 size={30} className="text-emerald-400" />
        </div>
        <h1 className="text-2xl font-bold mb-3">申请已提交</h1>
        <p className="text-gray-500 mb-8">
          管理团队将在 1-3 个工作日内审核你的入驻申请
        </p>
        <button className="btn-primary" onClick={() => router.push('/marketplace')}>
          返回集市
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto px-6 py-16">
      <button
        onClick={() => router.back()}
        className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-10"
      >
        <ArrowLeft size={16} />
        返回
      </button>

      <div className="mb-10">
        <span className="section-label flex items-center gap-2 mb-3">
          <Store size={13} />
          Merchant Application
        </span>
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight mb-3">入驻集市</h1>
        <p className="text-gray-400 leading-relaxed">
          成为 CampusHub 集市商家，发布你的知识商品、虚拟服务或校园周边。
          审核通过后即可上架商品，赚取 CHB 积分收益。
        </p>
      </div>

      <div className="glass-card p-8">
        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              店铺名称 <span className="text-red-400">*</span>
            </label>
            <input
              className="input"
              placeholder="例如：CampusHub 官方周边店"
              value={shopName}
              onChange={(e) => setShopName(e.target.value)}
              maxLength={128}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              店铺介绍
            </label>
            <textarea
              className="input"
              placeholder="介绍一下你的店铺和经营内容..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              maxLength={2000}
            />
          </div>

          <div className="flex items-start gap-3 p-4 rounded-xl bg-amber-500/5 border border-amber-500/15 text-sm text-amber-300/80">
            <AlertCircle size={16} className="shrink-0 mt-0.5" />
            <span>
              入驻后请遵守平台规则：商品信息真实有效，禁止虚假宣传。
              每笔交易将收取 10% 平台手续费。
            </span>
          </div>

          <button
            className="btn-primary w-full !py-4"
            disabled={submitting}
            onClick={handleSubmit}
          >
            <Send size={16} />
            {submitting ? '提交中...' : '提交申请'}
          </button>
        </div>
      </div>
    </div>
  );
}