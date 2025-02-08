import { useState } from "react";
import { Bars3Icon, BellIcon, MoonIcon, SunIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { Outlet } from "react-router";

const Layout: React.FC = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [isDarkMode, setIsDarkMode] = useState(false);

  return (
    <div className="flex min-h-screen bg-gray-200">
      {/* Sidebar */}
      <aside
        className={`${
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        } fixed inset-y-0 left-0 z-50 w-64 bg-violet-800 p-5 text-white transition-transform lg:relative lg:translate-x-0`}
      >
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold">Dashboard</h1>
          <button
            className="lg:hidden"
            onClick={() => setIsSidebarOpen(false)}
          >
            <XMarkIcon className="h-6 w-6" />
          </button>
        </div>
        <nav className="mt-5">
          <ul className="space-y-2">
            <li className="p-2 hover:bg-gray-700 rounded">Dashboard</li>
            <li className="p-2 hover:bg-gray-700 rounded">Settings</li>
            <li className="p-2 hover:bg-gray-700 rounded">Profile</li>
          </ul>
        </nav>
      </aside>

      {/* Main content */}
      <div className="flex flex-1 flex-col">
        {/* Navbar */}
        <header className="flex items-center justify-between bg-white p-4 shadow-md">
          <button className="lg:hidden" onClick={() => setIsSidebarOpen(true)}>
            <Bars3Icon className="h-6 w-6" />
          </button>
          <h2 className="ml-2 text-xl font-semibold">Dashboard</h2>
          <input
            type="text"
            placeholder="Search..."
            className="ml-2 rounded border px-3 py-1 text-gray-700 focus:outline-none"
          />
          <BellIcon className="ml-2 h-6 w-6 text-gray-700" />
          <button onClick={() => setIsDarkMode(!isDarkMode)} className="ml-2">
            {isDarkMode ? <SunIcon className="h-6 w-6" /> : <MoonIcon className="h-6 w-6 text-gray-700" />}
          </button>
          <div className="ml-2 h-8 w-8 rounded-full bg-gray-300"></div>
        </header>

        {/* Content */}
        <main className="p-4">
          <Outlet />
        </main>

        {/* Footer */}
        <footer className="mt-auto bg-white p-4 text-center shadow-md">
          &copy; 2024 Dashboard Inc.
        </footer>
      </div>
    </div>
  );
};

export default Layout;
