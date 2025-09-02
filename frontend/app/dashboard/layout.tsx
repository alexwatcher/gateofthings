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
      {/* Sidebar */}
      <aside className={`${sidebarOpen ? 'w-64' : 'w-16'} flex flex-col border-r border-green-500 bg-black/80 p-4 transition-all shadow-[0_0_15px_#00ff00]`}>
        <button onClick={() => setSidebarOpen(!sidebarOpen)} className="mb-6 flex items-center gap-2 rounded-md border border-green-500 p-2 hover:bg-green-500 hover:text-black">
          <Menu size={20} /> {sidebarOpen && <span>Menu</span>}
        </button>

        <nav className="flex flex-col gap-4">
          <button
            onClick={() => setActivePanel('home')}
            className={`flex items-center gap-2 rounded-md border border-green-500 p-2 w-full hover:bg-green-500 hover:text-black ${
              activePanel === 'home' ? 'bg-green-500 text-black' : ''
            }`}
          >
            <Home size={20} /> {sidebarOpen && <span>Home</span>}
          </button>

          <button
            onClick={() => setActivePanel('devices')}
            className={`flex items-center gap-2 rounded-md border border-green-500 p-2 w-full hover:bg-green-500 hover:text-black ${
              activePanel === 'devices' ? 'bg-green-500 text-black' : ''
            }`}
          >
            <Cpu size={20} /> {sidebarOpen && <span>Devices</span>}
          </button>

          <button
            onClick={() => setActivePanel('statistics')}
            className={`flex items-center gap-2 rounded-md border border-green-500 p-2 w-full hover:bg-green-500 hover:text-black ${
              activePanel === 'statistics' ? 'bg-green-500 text-black' : ''
            }`}
          >
            <BarChart2 size={20} /> {sidebarOpen && <span>Statistics</span>}
          </button>
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

        {/* Контент панели */}
        <main className="flex-1 p-6 overflow-auto">{renderPanel()}</main>
      </div>
    </div>
  );
}
