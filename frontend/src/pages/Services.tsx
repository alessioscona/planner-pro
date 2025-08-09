import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuthFetch } from '../useAuthFetch';

function useServicesApi() {
  const authFetch = useAuthFetch();
  return {
    fetchServices: async () => {
      const res = await authFetch('/api/v1/services');
      if (!res.ok) throw new Error('fetch services');
      return res.json();
    },
    createService: async (p: any) => {
      const res = await authFetch('/api/v1/services', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(p),
      });
      if (!res.ok) throw new Error('create');
      return res.json();
    },
  };
}

export default function ServicesPage(){
  const qc = useQueryClient();
  const { fetchServices, createService } = useServicesApi();
  const { data } = useQuery({
    queryKey: ['services'],
    queryFn: fetchServices,
  });
  const mutation = useMutation({
    mutationFn: createService,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['services'] });
    },
  });
  const [name, setName] = useState('');
  const [duration, setDuration] = useState(30);
  const [price, setPrice] = useState(5000);

  return (
    <div>
      <h2>Services</h2>
      <form onSubmit={e => {
        e.preventDefault();
        mutation.mutate({ name, duration_minutes: Number(duration), price_cents: Number(price) });
        setName('');
      }}>
        <input value={name} onChange={e => setName(e.target.value)} placeholder="Name" required />
        <input type="number" value={duration} onChange={e => setDuration(Number(e.target.value))} />
        <input type="number" value={price} onChange={e => setPrice(Number(e.target.value))} />
        <button type="submit">Add</button>
      </form>
      <ul>
        {data && data.map((s: any) => <li key={s.id}>{s.name} - {s.duration_minutes}m</li>)}
      </ul>
    </div>
  );
}
