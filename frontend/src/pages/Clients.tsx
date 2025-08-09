import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuthFetch } from '../useAuthFetch';

import { Table, Button, Modal, Form, Input, Space, message } from 'antd';

function useClientsApi() {
  const authFetch = useAuthFetch();
  return {
    fetchClients: async () => {
      const res = await authFetch('/api/v1/clients');
      if (!res.ok) throw new Error('Errore nel recupero clienti');
      return res.json();
    },
    createClient: async (payload: any) => {
      const res = await authFetch('/api/v1/clients', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      if (!res.ok) throw new Error('Errore nella creazione cliente');
      return res.json();
    },
  };
}

export default function ClientsPage(){
  const qc = useQueryClient();
  const { fetchClients, createClient } = useClientsApi();
  const { data, isLoading, error } = useQuery({
    queryKey: ['clients'],
    queryFn: fetchClients,
  });
  const mutation = useMutation({
    mutationFn: createClient,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['clients'] });
    },
  });

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const columns = [
    {
      title: 'Nome',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Telefono',
      dataIndex: 'phone',
      key: 'phone',
    },
  ];

  return (
    <div style={{ padding: 20 }}>
      <h2>Clients</h2>
      <Button type="primary" onClick={() => setIsModalOpen(true)} style={{ marginBottom: 16 }}>
        Aggiungi Cliente
      </Button>

      <Table
        loading={isLoading}
        dataSource={data || []}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 5 }}
        locale={{ emptyText: error ? 'Errore nel caricamento dati' : 'Nessun cliente trovato' }}
      />

      <Modal
        title="Aggiungi Nuovo Cliente"
        open={isModalOpen}
        onCancel={() => {
          setIsModalOpen(false);
          form.resetFields();
        }}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={(values) => mutation.mutate(values)}
        >
          <Form.Item
            label="Nome"
            name="name"
            rules={[{ required: true, message: 'Inserisci il nome!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label="Email"
            name="email"
            rules={[
              { required: true, message: 'Inserisci l\'email!' },
              { type: 'email', message: 'Inserisci un\'email valida!' },
            ]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label="Telefono"
            name="phone"
            rules={[{ required: true, message: 'Inserisci il numero di telefono!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setIsModalOpen(false);
                form.resetFields();
              }}>
                Annulla
              </Button>
              <Button type="primary" htmlType="submit" loading={mutation.isLoading}>
                Aggiungi
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
