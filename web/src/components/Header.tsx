import { useState } from "react";
import { MagnifyingGlassIcon, Bars3Icon } from "@heroicons/react/24/solid";

const Header: React.FC = () => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const toggleDropdown = () => {
    setIsDropdownOpen((prev) => !prev);
  };

  return (
    <header className="shadow-bottom sticky top-0 z-20 lg:ml-64 h-16 bg-white">
      <div className="container mx-auto flex h-full items-center justify-between px-6 text-purple-600">
        <button
          className="focus:shadow-outline-purple -ml-1 mr-5 rounded-md p-1 focus:outline-none lg:hidden"
          aria-label="Menu"
        >
          <Bars3Icon className="h-4 w-4" aria-hidden="true" />
        </button>

        <div className="flex flex-1 justify-center">
          <div className="relative mr-6 w-full max-w-xl focus-within:text-purple-500">
            <div className="absolute inset-y-0 flex items-center pl-2">
              <MagnifyingGlassIcon className="h-4 w-4" aria-hidden="true" />
            </div>
            <input
              className="text w-full rounded-lg bg-gray-100 py-2 pl-8 text-gray-700 outline-purple-500"
              placeholder="Search something..."
              aria-label="Search"
            />
          </div>
        </div>

        <ul className="flex flex-shrink-0 items-center space-x-6">
          <li className="relative">
            <button
              onClick={toggleDropdown}
              className="focus:shadow-outline-purple rounded-md focus:outline-none"
              aria-label="Toggle dropdown"
            >
              <img
                className="h-10 w-10 rounded-full border-2 border-purple-400 object-cover"
                src="https://picsum.photos/150"
                alt="User avatar"
              />
            </button>
            {isDropdownOpen && (
              <ul className="absolute right-0 mt-2 w-48 rounded-lg bg-gray-50 text-gray-700 shadow-lg">
                <li className="cursor-pointer p-2 text-sm hover:bg-gray-100">
                  Profile
                </li>
                <li className="cursor-pointer p-2 text-sm hover:bg-gray-100">
                  Settings
                </li>
                <li className="cursor-pointer p-2 text-sm hover:bg-gray-100">
                  Logout
                </li>
              </ul>
            )}
          </li>
        </ul>
      </div>
    </header>
  );
};

export default Header;
