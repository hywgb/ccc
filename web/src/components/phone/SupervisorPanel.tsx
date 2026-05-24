import { useState, useEffect } from 'react';
import { Card, Table, Tag, Button, Space, Select, message, Badge } from 'antd';
import { EyeOutlined, AudioOutlined, ThunderboltOutlined, StopOutlined, ReloadOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { callControlApi } from '../../api/endpoints';
import api from '../../api/client';

interface ActiveCall {
  id: number;
  agent_name: string;
  agent_id: number;
  customer_phone: string;
  direction: string;
  duration: number;
  status: string;
  skill_group: string;
}

export default function SupervisorPanel() {
  const [calls, setCalls] = useState<ActiveCall[]>([]);
  const [loading, setLoading] = useState(false);
  const [filterSkill, setFilterSkill] = useState<string>('');

  const loadCalls = async () => {
    setLoading(true);
    try {
      const res = await api.get('/supervisor/active-calls');
      setCalls(Array.isArray(res.data) ? res.data : res.data?.items || []);
    } catch { /* ignore */ }
    setLoading(false);
  };

  useEffect(() => { loadCalls(); const t = setInterval(loadCalls, 5000); return () => clearInterval(t); }, []);

  const handleAction = async (callId: number, action: string) => {
    try {
      switch (action) {
        case 'monitor': await callControlApi.monitor(callId); break;
        case 'whisper': await callControlApi.whisper(callId); break;
        case 'barge':   await callControlApi.barge(callId); break;
        case 'intercept': await callControlApi.intercept(callId); break;
        case 'coach':   await callControlApi.coach(callId); break;
      }
      message.success(`${action} 操作成功`);
    } catch {
      message.error('操作失败');
    }
  };

  const columns: ColumnsType<ActiveCall> = [
    { title: '坐席', dataIndex: 'agent_name', width: 100 },
    { title: '客户号码', dataIndex: 'customer_phone', width: 130 },
    {
      title: '方向', dataIndex: 'direction', width: 70,
      render: (v) => <Tag color={v === 'inbound' ? 'blue' : 'green'}>{v === 'inbound' ? '呼入' : '呼出'}</Tag>,
    },
    { title: '技能组', dataIndex: 'skill_group', width: 100 },
    {
      title: '时长', dataIndex: 'duration', width: 80,
      render: (s: number) => `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`,
    },
    {
      title: '状态', dataIndex: 'status', width: 80,
      render: (v) => <Badge status={v === 'talking' ? 'processing' : 'default'} text={v === 'talking' ? '通话中' : v} />,
    },
    {
      title: '操作', key: 'actions', width: 280,
      render: (_, record) => (
        <Space size={4}>
          <Button size="small" icon={<EyeOutlined />} onClick={() => handleAction(record.id, 'monitor')}>监听</Button>
          <Button size="small" icon={<AudioOutlined />} onClick={() => handleAction(record.id, 'whisper')}>耳语</Button>
          <Button size="small" icon={<ThunderboltOutlined />} onClick={() => handleAction(record.id, 'barge')} danger>强插</Button>
          <Button size="small" icon={<StopOutlined />} onClick={() => handleAction(record.id, 'intercept')} danger>强拆</Button>
        </Space>
      ),
    },
  ];

  const filtered = filterSkill ? calls.filter((c) => c.skill_group === filterSkill) : calls;
  const skillOptions = [...new Set(calls.map((c) => c.skill_group))].filter(Boolean).map((s) => ({ value: s, label: s }));

  return (
    <Card
      title="实时监控"
      extra={
        <Space>
          <Select allowClear placeholder="筛选技能组" options={skillOptions} value={filterSkill || undefined} onChange={(v) => setFilterSkill(v || '')} style={{ width: 150 }} />
          <Button icon={<ReloadOutlined />} onClick={loadCalls}>刷新</Button>
        </Space>
      }
    >
      <Table<ActiveCall> columns={columns} dataSource={filtered} rowKey="id" loading={loading} size="small" pagination={false} />
    </Card>
  );
}
