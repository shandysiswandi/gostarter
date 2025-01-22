import  { useState } from "react";

const Sidebar: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="bg-violet-800 text-white lg:w-64">
      <button className="p-4 lg:hidden" onClick={() => setIsOpen(!isOpen)}>
        {isOpen ? "Close" : "Open"} Menu
      </button>
      <div className={`lg:block ${isOpen ? "block" : "hidden"}`}>
        <a href="#" className="block rounded p-2 hover:bg-gray-700">
          Dashboard
        </a>
        <a href="#" className="block rounded p-2 hover:bg-gray-700">
          Profile
        </a>
        <a href="#" className="block rounded p-2 hover:bg-gray-700">
          Settings
        </a>
        <a href="#" className="block rounded p-2 hover:bg-gray-700">
          Logout
        </a>
      </div>
    </div>
  );
};

export default Sidebar;
