import { useState } from 'react';
import { Card, Tabs, Table, Button, Form, Input, Select, Upload, Space, Tag, message, Modal, Progress, Empty } from 'antd';
import {
  AudioOutlined, BarChartOutlined, ExperimentOutlined, UploadOutlined,
  PlayCircleOutlined, PlusOutlined, ReloadOutlined, RobotOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import api from '../../api/client';

// --- Voice Cloning ---
interface VoiceCloneTask {
  id: number;
  name: string;
  status: 'pending' | 'training' | 'ready' | 'failed';
  progress: number;
  created_at: string;
}

function VoiceCloningTab() {
  const [tasks, setTasks] = useState<VoiceCloneTask[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  const load = async () => {
    setLoading(true);
    try {
      const res = await api.get('/ai/voice-clone/tasks');
      setTasks(Array.isArray(res.data) ? res.data : res.data?.items || []);
    } catch { /* ignore */ }
    setLoading(false);
  };

  const handleCreate = async () => {
    const values = await form.validateFields();
    try {
      await api.post('/ai/voice-clone/tasks', values);
      message.success('训练任务已创建');
      setModalOpen(false);
      form.resetFields();
      load();
    } catch {
      message.error('创建失败');
    }
  };

  const statusColor: Record<string, string> = { pending: 'default', training: 'processing', ready: 'success', failed: 'error' };
  const statusLabel: Record<string, string> = { pending: '等待中', training: '训练中', ready: '就绪', failed: '失败' };

  const columns: ColumnsType<VoiceCloneTask> = [
    { title: '名称', dataIndex: 'name' },
    {
      title: '状态', dataIndex: 'status', width: 100,
      render: (v) => <Tag color={statusColor[v]}>{statusLabel[v] || v}</Tag>,
    },
    {
      title: '进度', dataIndex: 'progress', width: 200,
      render: (v, record) => record.status === 'training' ? <Progress percent={v} size="small" /> : '-',
    },
    { title: '创建时间', dataIndex: 'created_at', width: 160 },
    {
      title: '操作', key: 'action', width: 120,
      render: (_, record) => record.status === 'ready' ? (
        <Button size="small" icon={<PlayCircleOutlined />} onClick={() => message.info('试听功能开发中')}>试听</Button>
      ) : null,
    },
  ];

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: 12 }}>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>新建训练</Button>
        </Space>
      </div>
      <Table<VoiceCloneTask> columns={columns} dataSource={tasks} rowKey="id" loading={loading} size="small" pagination={false} />
      <Modal title="新建声纹复刻任务" open={modalOpen} onOk={handleCreate} onCancel={() => { setModalOpen(false); form.resetFields(); }} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}><Input placeholder="声纹名称" /></Form.Item>
          <Form.Item name="audio_url" label="训练音频"><Upload maxCount={1}><Button icon={<UploadOutlined />}>上传音频</Button></Upload></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={2} /></Form.Item>
        </Form>
      </Modal>
    </>
  );
}

// --- Conversation Analytics ---
function ConversationAnalyticsTab() {
  const [result, setResult] = useState<Record<string, unknown> | null>(null);
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  const runAnalysis = async () => {
    const values = await form.validateFields();
    setLoading(true);
    try {
      const res = await api.post('/ai/conversation-analytics/analyze', values);
      setResult(res.data);
    } catch {
      message.error('分析失败');
    }
    setLoading(false);
  };

  return (
    <div>
      <Card title="对话分析" size="small">
        <Form form={form} layout="vertical">
          <Form.Item name="type" label="分析类型" rules={[{ required: true }]}>
            <Select options={[
              { value: 'intent_mining', label: '意图挖掘' },
              { value: 'sop_discovery', label: 'SOP 发现' },
              { value: 'sales_script', label: '金牌话术提取' },
              { value: 'topic_clustering', label: '话题聚类' },
            ]} />
          </Form.Item>
          <Form.Item name="date_range" label="时间范围">
            <Select options={[
              { value: '7d', label: '近 7 天' },
              { value: '30d', label: '近 30 天' },
              { value: '90d', label: '近 90 天' },
            ]} />
          </Form.Item>
          <Form.Item name="skill_group_id" label="技能组(可选)"><Input placeholder="技能组ID" /></Form.Item>
          <Button type="primary" onClick={runAnalysis} loading={loading}>开始分析</Button>
        </Form>
      </Card>
      {result && (
        <Card title="分析结果" size="small" style={{ marginTop: 16 }}>
          <pre style={{ maxHeight: 400, overflow: 'auto', background: '#f5f5f5', padding: 12, borderRadius: 4, fontSize: 13 }}>
            {JSON.stringify(result, null, 2)}
          </pre>
        </Card>
      )}
    </div>
  );
}

// --- Training ---
function TrainingTab() {
  const [questions, setQuestions] = useState<{ id: number; question: string; answer: string; score?: number }[]>([]);
  const [loading, setLoading] = useState(false);

  const generateQuestions = async () => {
    setLoading(true);
    try {
      const res = await api.post('/ai/training/generate-questions', { count: 10 });
      setQuestions(Array.isArray(res.data) ? res.data : res.data?.items || []);
    } catch {
      message.error('生成失败');
    }
    setLoading(false);
  };

  const evaluateCall = async () => {
    try {
      const res = await api.post('/ai/training/evaluate', { call_id: 0 });
      message.success(`评估完成: ${JSON.stringify(res.data?.score || res.data)}`);
    } catch {
      message.error('评估失败');
    }
  };

  return (
    <>
      <Space style={{ marginBottom: 12 }}>
        <Button type="primary" icon={<ExperimentOutlined />} onClick={generateQuestions} loading={loading}>生成考题</Button>
        <Button icon={<BarChartOutlined />} onClick={evaluateCall}>模拟通话评估</Button>
      </Space>
      {questions.length > 0 ? (
        <Table
          size="small"
          rowKey="id"
          dataSource={questions}
          columns={[
            { title: '#', dataIndex: 'id', width: 50 },
            { title: '问题', dataIndex: 'question' },
            { title: '参考答案', dataIndex: 'answer', ellipsis: true },
            { title: '得分', dataIndex: 'score', width: 80, render: (v) => v ? <Tag color="blue">{v}</Tag> : '-' },
          ]}
          pagination={false}
        />
      ) : (
        <Empty description={'点击"生成考题"开始'} image={Empty.PRESENTED_IMAGE_SIMPLE} />
      )}
    </>
  );
}

// --- Main Page ---
export default function AdvancedAiPage() {
  return (
    <>
      <h2><RobotOutlined /> 高级 AI</h2>
      <Tabs items={[
        { key: 'voice', label: <Space><AudioOutlined />声纹复刻</Space>, children: <VoiceCloningTab /> },
        { key: 'analytics', label: <Space><BarChartOutlined />对话分析</Space>, children: <ConversationAnalyticsTab /> },
        { key: 'training', label: <Space><ExperimentOutlined />智能培训</Space>, children: <TrainingTab /> },
      ]} />
    </>
  );
}
