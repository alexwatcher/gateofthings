'use client';
import { useState } from 'react';
import { Menu, User, Home, Cpu, BarChart2 } from 'lucide-react';
import HomePanel from './components/homepanel';
import DevicesPanel from './components/devicespanel';
import StatisticsPanel from './components/statisticspanel';

export default function DashboardLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [activePanel, setActivePanel] = useState<'home' | 'devices' | 'statistics'>('home');
  const [userMenuOpen, setUserMenuOpen] = useState(false);

  const renderPanel = () => {
    switch (activePanel) {
      case 'home':
        return <HomePanel />;
      case 'devices':
        return <DevicesPanel />;
      case 'statistics':
        return <StatisticsPanel />;
    }
  };

  return (
    <div className="flex h-screen bg-black text-green-400 font-mono">
<aside
  className={`flex flex-col border-r border-green-500 bg-black/80 p-4 transition-all shadow-[0_0_15px_#00ff00] ${
    sidebarOpen ? 'w-64' : 'w-16'
  }`}
>
  {/* Menu Button */}
<div className="flex justify-end mb-6 flex-none">
  <button
    onClick={() => setSidebarOpen(!sidebarOpen)}
    className="flex h-10 w-10 items-center justify-center rounded-md border border-green-500 bg-black text-green-400 transition-colors hover:bg-green-500 hover:text-black"
  >
    {sidebarOpen ? (
      <span className="text-lg font-bold">×</span>
    ) : (
      <div className="flex flex-col justify-between h-4 w-5">
        <span className="block h-0.5 w-full bg-green-400"></span>
        <span className="block h-0.5 w-full bg-green-400"></span>
        <span className="block h-0.5 w-full bg-green-400"></span>
      </div>
    )}
  </button>
</div>

  {/* Navigation */}
  <nav className="flex flex-col gap-4">
    {[
      { name: 'home', icon: <Home size={20} /> },
      { name: 'devices', icon: <Cpu size={20} /> },
      { name: 'statistics', icon: <BarChart2 size={20} /> },
    ].map((item) => (
      <button
        key={item.name}
        onClick={() => setActivePanel(item.name as any)}
        className={`flex items-center justify-center gap-2 rounded-md border border-green-500 p-2 w-full h-10 hover:bg-green-500 hover:text-black ${
          activePanel === item.name ? 'bg-green-500 text-black' : ''
        } ${sidebarOpen ? 'justify-start px-2' : ''}`}
      >
        {item.icon} {sidebarOpen && <span>{item.name.charAt(0).toUpperCase() + item.name.slice(1)}</span>}
      </button>
    ))}
  </nav>
</aside>

      {/* Main */}
      <div className="flex flex-1 flex-col">
        {/* Topbar */}
        <header className="flex items-center justify-between border-b border-green-500 bg-black/80 px-6 py-4 shadow-[0_0_15px_#00ff00]">
          <div>&nbsp;</div>
          <div className="relative">
            <button onClick={() => setUserMenuOpen(!userMenuOpen)} className="flex items-center gap-2 rounded-md border border-green-500 p-2 hover:bg-green-500 hover:text-black">
              <User size={20} /> <span>neo</span>
            </button>

            {userMenuOpen && (
              <div className="absolute right-0 mt-2 w-40 rounded-md border border-green-500 bg-black/90 shadow-[0_0_10px_#00ff00] flex flex-col">
                <button className="p-2 text-green-400 hover:bg-green-500 hover:text-black">Settings</button>
                <button className="p-2 text-green-400 hover:bg-green-500 hover:text-black">Logout</button>
              </div>
            )}
          </div>
        </header>

        {/* Penel content */}
        <main className="flex-1 p-6 overflow-auto">{renderPanel()}</main>
      </div>
    </div>
  );
}
