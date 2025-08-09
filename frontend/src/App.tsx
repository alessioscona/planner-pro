import React from 'react'
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import { Layout, Menu, Typography } from 'antd'
import ClientsPage from './pages/Clients'
import ServicesPage from './pages/Services'
import AppointmentsPage from './pages/Appointments'
import { RequireAuth } from './RequireAuth'

const { Header, Content } = Layout
const { Title } = Typography

export default function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Header style={{ display: 'flex', alignItems: 'center' }}>
          <Title level={3} style={{ color: 'white', margin: 0, flex: '1 0 auto' }}>
            Scheduler
          </Title>
          <Menu
            theme="dark"
            mode="horizontal"
            selectable={false}
            items={[
              { key: 'home', label: <Link to="/">Home</Link> },
              { key: 'clients', label: <Link to="/clients">Clients</Link> },
              { key: 'services', label: <Link to="/services">Services</Link> },
              { key: 'appointments', label: <Link to="/appointments">Appointments</Link> },
            ]}
          />
        </Header>
        <Content style={{ padding: '24px', maxWidth: 1200, margin: '0 auto' }}>
          <RequireAuth>
            <Routes>
              <Route path="/" element={<div>Welcome to Scheduler</div>} />
              <Route path="/clients" element={<ClientsPage />} />
              <Route path="/services" element={<ServicesPage />} />
              <Route path="/appointments" element={<AppointmentsPage />} />
            </Routes>
          </RequireAuth>
        </Content>
      </Layout>
    </BrowserRouter>
  )
}
