import { useState } from 'react';
import { Modal, Tabs, Select, Button, Space, message } from 'antd';
import { callControlApi, skillGroupApi, agentApi } from '../../api/endpoints';

interface TransferModalProps {
  open: boolean;
  callId: number | null;
  onClose: () => void;
}

export default function TransferModal({ open, callId, onClose }: TransferModalProps) {
  const [tab, setTab] = useState('skill');
  const [target, setTarget] = useState<string>('');
  const [mode, setMode] = useState<'blind' | 'attended'>('blind');
  const [skills, setSkills] = useState<{ value: string; label: string }[]>([]);
  const [agents, setAgents] = useState<{ value: string; label: string }[]>([]);

  const loadOptions = async () => {
    try {
      const [skillRes, agentRes] = await Promise.all([skillGroupApi.list(), agentApi.list()]);
      const skillData = Array.isArray(skillRes.data) ? skillRes.data : skillRes.data?.items || [];
      const agentData = Array.isArray(agentRes.data) ? agentRes.data : agentRes.data?.items || [];
      setSkills(skillData.map((s: Record<string, unknown>) => ({ value: String(s.id), label: String(s.name) })));
      setAgents(agentData.map((a: Record<string, unknown>) => ({ value: String(a.id), label: String(a.display_name || a.name) })));
    } catch { /* ignore */ }
  };

  const handleTransfer = async () => {
    if (!callId || !target) { message.warning('请选择转接目标'); return; }
    try {
      const data = tab === 'skill'
        ? { skill_group_id: Number(target) }
        : tab === 'agent'
          ? { agent_id: Number(target) }
          : { external_number: target };

      if (mode === 'blind') {
        await callControlApi.blindTransfer(callId, data);
      } else {
        await callControlApi.attendedTransfer(callId, data);
      }
      message.success('转接成功');
      onClose();
    } catch {
      message.error('转接失败');
    }
  };

  return (
    <Modal title="转接通话" open={open} onCancel={onClose} footer={null} afterOpenChange={(v) => v && loadOptions()}>
      <Tabs activeKey={tab} onChange={(k) => { setTab(k); setTarget(''); }} items={[
        { key: 'skill', label: '技能组' },
        { key: 'agent', label: '坐席' },
        { key: 'external', label: '外线' },
      ]} />

      <div style={{ marginBottom: 16 }}>
        {tab === 'skill' && (
          <Select placeholder="选择技能组" options={skills} value={target || undefined} onChange={setTarget} style={{ width: '100%' }} showSearch optionFilterProp="label" />
        )}
        {tab === 'agent' && (
          <Select placeholder="选择坐席" options={agents} value={target || undefined} onChange={setTarget} style={{ width: '100%' }} showSearch optionFilterProp="label" />
        )}
        {tab === 'external' && (
          <Select mode="tags" placeholder="输入外线号码" value={target ? [target] : []} onChange={(v) => setTarget(v[0] || '')} style={{ width: '100%' }} />
        )}
      </div>

      <Space style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Select value={mode} onChange={setMode} options={[
          { value: 'blind', label: '盲转' },
          { value: 'attended', label: '咨询转' },
        ]} style={{ width: 100 }} />
        <Button type="primary" onClick={handleTransfer}>确定转接</Button>
      </Space>
    </Modal>
  );
}
