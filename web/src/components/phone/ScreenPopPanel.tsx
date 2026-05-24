import { useState, useEffect } from 'react';
import { Card, Descriptions, Tag, Tabs, Spin, Empty, Space } from 'antd';
import { UserOutlined, PhoneOutlined, HistoryOutlined } from '@ant-design/icons';
import api from '../../api/client';

interface CustomerInfo {
  name: string;
  phone: string;
  company: string;
  level: string;
  last_contact: string;
  notes: string;
  tags: string[];
  history: { id: number; date: string; type: string; summary: string }[];
}

export default function ScreenPopPanel({ callerNumber }: { callerNumber?: string }) {
  const [customer, setCustomer] = useState<CustomerInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [iframeUrl, setIframeUrl] = useState<string>('');

  useEffect(() => {
    if (!callerNumber) return;
    setLoading(true);
    api.get('/screen-pop/lookup', { params: { phone: callerNumber } })
      .then((res) => {
        setCustomer(res.data?.customer || null);
        setIframeUrl(res.data?.iframe_url || '');
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [callerNumber]);

  if (loading) return <Card size="small"><Spin style={{ display: 'block', padding: 40 }} /></Card>;
  if (!callerNumber) return null;

  return (
    <Card title={<Space><UserOutlined /> 来电弹屏</Space>} size="small">
      <Tabs size="small" items={[
        {
          key: 'info',
          label: <Space><UserOutlined />客户信息</Space>,
          children: customer ? (
            <>
              <Descriptions size="small" column={2} bordered>
                <Descriptions.Item label="姓名">{customer.name}</Descriptions.Item>
                <Descriptions.Item label="电话">{customer.phone}</Descriptions.Item>
                <Descriptions.Item label="公司">{customer.company}</Descriptions.Item>
                <Descriptions.Item label="等级"><Tag color="blue">{customer.level}</Tag></Descriptions.Item>
                <Descriptions.Item label="最近联系">{customer.last_contact}</Descriptions.Item>
                <Descriptions.Item label="标签">
                  {customer.tags?.map((t) => <Tag key={t}>{t}</Tag>)}
                </Descriptions.Item>
              </Descriptions>
              {customer.notes && (
                <div style={{ marginTop: 8, padding: 8, background: '#fafafa', borderRadius: 4, fontSize: 13 }}>
                  {customer.notes}
                </div>
              )}
            </>
          ) : (
            <Empty description="未找到客户信息" image={Empty.PRESENTED_IMAGE_SIMPLE} />
          ),
        },
        {
          key: 'history',
          label: <Space><HistoryOutlined />历史记录</Space>,
          children: customer?.history?.length ? (
            <div style={{ maxHeight: 300, overflowY: 'auto' }}>
              {customer.history.map((h) => (
                <div key={h.id} style={{ padding: '6px 0', borderBottom: '1px solid #f0f0f0' }}>
                  <Space>
                    <Tag>{h.type}</Tag>
                    <span style={{ fontSize: 12, color: '#999' }}>{h.date}</span>
                  </Space>
                  <div style={{ fontSize: 13, marginTop: 4 }}>{h.summary}</div>
                </div>
              ))}
            </div>
          ) : (
            <Empty description="暂无历史" image={Empty.PRESENTED_IMAGE_SIMPLE} />
          ),
        },
        ...(iframeUrl ? [{
          key: 'iframe',
          label: <Space><PhoneOutlined />业务系统</Space>,
          children: <iframe src={iframeUrl} style={{ width: '100%', height: 400, border: 'none' }} title="弹屏" />,
        }] : []),
      ]} />
    </Card>
  );
}
