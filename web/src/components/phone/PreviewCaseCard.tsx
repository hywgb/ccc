import { useState, useEffect, useRef } from 'react';
import { Card, Descriptions, Button, Space, Tag, Spin, message } from 'antd';
import { PhoneOutlined, ForwardOutlined, ClockCircleOutlined } from '@ant-design/icons';
import api from '../../api/client';

interface PreviewCase {
  id: number;
  campaign_id: number;
  campaign_name: string;
  customer_name: string;
  phone: string;
  company: string;
  priority: number;
  attempt_count: number;
  custom_fields: Record<string, string>;
  notes: string;
  expires_in: number;
}

export default function PreviewCaseCard() {
  const [caseData, setCaseData] = useState<PreviewCase | null>(null);
  const [loading, setLoading] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const wasCountingRef = useRef(false);

  const loadCase = async () => {
    setLoading(true);
    try {
      const res = await api.get('/campaigns/preview/current');
      setCaseData(res.data || null);
      if (res.data?.expires_in) setCountdown(res.data.expires_in);
    } catch {
      setCaseData(null);
    }
    setLoading(false);
  };

  useEffect(() => { loadCase(); }, []);

  useEffect(() => {
    if (countdown <= 0) return;
    wasCountingRef.current = true;
    const t = setInterval(() => {
      setCountdown((c) => (c <= 1 ? 0 : c - 1));
    }, 1000);
    return () => clearInterval(t);
  }, [countdown > 0]);

  useEffect(() => {
    if (countdown === 0 && wasCountingRef.current) {
      wasCountingRef.current = false;
      handleSkip();
    }
  }, [countdown]);

  const handleDial = async () => {
    if (!caseData) return;
    try {
      await api.post(`/campaigns/${caseData.campaign_id}/cases/${caseData.id}/dial`);
      message.success('已发起呼叫');
      setCaseData(null);
    } catch {
      message.error('拨号失败');
    }
  };

  const handleSkip = async () => {
    if (!caseData) return;
    try {
      await api.post(`/campaigns/${caseData.campaign_id}/cases/${caseData.id}/skip`);
      loadCase();
    } catch { /* ignore */ }
  };

  if (loading) return <Card size="small"><Spin style={{ display: 'block', padding: 24 }} /></Card>;
  if (!caseData) return null;

  return (
    <Card
      title={<Space><PhoneOutlined /> Preview 预览案例</Space>}
      size="small"
      extra={<Tag icon={<ClockCircleOutlined />} color={countdown < 10 ? 'red' : 'orange'}>{countdown}s</Tag>}
    >
      <Descriptions size="small" column={2} bordered>
        <Descriptions.Item label="客户">{caseData.customer_name}</Descriptions.Item>
        <Descriptions.Item label="电话">{caseData.phone}</Descriptions.Item>
        <Descriptions.Item label="公司">{caseData.company}</Descriptions.Item>
        <Descriptions.Item label="优先级"><Tag color="blue">{caseData.priority}</Tag></Descriptions.Item>
        <Descriptions.Item label="活动">{caseData.campaign_name}</Descriptions.Item>
        <Descriptions.Item label="尝试次数">{caseData.attempt_count}</Descriptions.Item>
        {Object.entries(caseData.custom_fields || {}).map(([k, v]) => (
          <Descriptions.Item key={k} label={k}>{v}</Descriptions.Item>
        ))}
      </Descriptions>
      {caseData.notes && (
        <div style={{ marginTop: 8, padding: 8, background: '#fffbe6', borderRadius: 4, fontSize: 13 }}>
          {caseData.notes}
        </div>
      )}
      <Space style={{ marginTop: 12, display: 'flex', justifyContent: 'flex-end' }}>
        <Button icon={<ForwardOutlined />} onClick={handleSkip}>跳过</Button>
        <Button type="primary" icon={<PhoneOutlined />} onClick={handleDial} style={{ background: '#52c41a' }}>拨号</Button>
      </Space>
    </Card>
  );
}
