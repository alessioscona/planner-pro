import React, { useState } from 'react'
import { useAuthFetch } from '../useAuthFetch';

export default function Appointments(){
  const [clientId, setClientId] = useState('')
  const [serviceId, setServiceId] = useState('')
  const [start, setStart] = useState('')

  const authFetch = useAuthFetch();
  const submit = async (e:any)=>{
    e.preventDefault();
    const res = await authFetch('/api/v1/appointments', {
      method:'POST',
      headers:{'Content-Type':'application/json'},
      body: JSON.stringify({client_id: clientId, service_id: serviceId, start_at: start})
    });
    if (!res.ok) { alert('err:'+res.status); return }
    alert('created');
  }

  return (
    <div>
      <h2>Appointments</h2>
      <form onSubmit={submit}>
        <input placeholder="client uuid" value={clientId} onChange={e=>setClientId(e.target.value)} />
        <input placeholder="service uuid" value={serviceId} onChange={e=>setServiceId(e.target.value)} />
        <input placeholder="start (RFC3339)" value={start} onChange={e=>setStart(e.target.value)} />
        <button type="submit">Schedule</button>
      </form>
    </div>
  )
}
