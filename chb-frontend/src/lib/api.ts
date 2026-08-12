// API 基础配置
export const API_BASE = process.env.NEXT_PUBLIC_API_BASE || '';

// 用户身份（开发/测试环境通过 X-User-ID 模拟）
export function getUserID(): number {
  if (typeof window === 'undefined') return 0;
  const stored = localStorage.getItem('chb_user_id');
  if (stored) return Number(stored);
  return 1;
}

export function setUserID(id: number) {
  localStorage.setItem('chb_user_id', String(id));
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
}

export interface ApiResponse<T = Record<string, unknown>> {
  code: number;
  message: string;
  data: T;
}

export interface Paginated<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

export async function apiRequest<T = Record<string, unknown>>(
  path: string,
  options: RequestOptions = {}
): Promise<ApiResponse<T>> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    Accept: 'application/json, text/plain, */*',
    ...options.headers,
  };

  // 开发环境模拟用户
  headers['X-User-ID'] = String(getUserID());

  const res = await fetch(`${API_BASE}${path}`, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  return res.json() as Promise<ApiResponse<T>>;
}

// 格式化 CHB 数量
export function formatCHB(amount: number): string {
  return new Intl.NumberFormat('zh-CN').format(amount) + ' CHB';
}

// 格式化日期
export function formatDate(dateStr: string): string {
  const d = new Date(dateStr);
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' });
}

// 生成幂等键
export function genIdempotencyKey(prefix: string): string {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`;
}