import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuthFetch } from '../useAuthFetch';

function useClientsApi() {
  const authFetch = useAuthFetch();
  return {
    fetchClients: async () => {
      const res = await authFetch('/api/v1/clients');
      if (!res.ok) throw new Error('fetch clients');
      return res.json();
    },
    createClient: async (payload: any) => {
      const res = await authFetch('/api/v1/clients', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      if (!res.ok) throw new Error('create');
      return res.json();
    },
  };
}

export default function ClientsPage(){
  const qc = useQueryClient();
  const { fetchClients, createClient } = useClientsApi();
  const { data } = useQuery({
    queryKey: ['clients'],
    queryFn: fetchClients,
  });
  const mutation = useMutation({
    mutationFn: createClient,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['clients'] });
    },
  });
  const [name, setName] = useState('');

  return (
    <div>
      <h2>Clients</h2>
      <form onSubmit={e => {
        e.preventDefault();
        mutation.mutate({ name });
        setName('');
      }}>
        <input value={name} onChange={e => setName(e.target.value)} placeholder="Name" required />
        <button type="submit">Add</button>
      </form>
      <ul>
        {data && data.map((c: any) => <li key={c.id}>{c.name}</li>)}
      </ul>
    </div>
  );
}
