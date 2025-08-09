
import { RequireAuth } from './RequireAuth';
import React from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import ClientsPage from './pages/Clients';
import ServicesPage from './pages/Services';
import AppointmentsPage from './pages/Appointments';

export default function App() {
  return (
    <BrowserRouter>
      <header style={{ padding: 10, borderBottom: '1px solid #eee' }}>
        <Link to="/">Home</Link> | <Link to="/clients">Clients</Link> | <Link to="/services">Services</Link> | <Link to="/appointments">Appointments</Link>
      </header>
      <main style={{ padding: 10 }}>
        <RequireAuth>
          <Routes>
            <Route path="/" element={<div>Welcome to Scheduler</div>} />
            <Route path="/clients" element={<ClientsPage />} />
            <Route path="/services" element={<ServicesPage />} />
            <Route path="/appointments" element={<AppointmentsPage />} />
          </Routes>
        </RequireAuth>
      </main>
    </BrowserRouter>
  );
}
