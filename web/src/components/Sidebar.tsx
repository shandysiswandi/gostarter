import { PATH_DASHBOARD } from "@/routes";
import {
  HomeIcon,
  InboxIcon,
  BellIcon,
  ChatBubbleBottomCenterTextIcon,
  DocumentTextIcon,
  UsersIcon,
  ArrowRightStartOnRectangleIcon,
  Cog8ToothIcon,
} from "@heroicons/react/24/outline";
import { Link } from "react-router";
// https://www.creative-tim.com/twcomponents/component/sidebar-navigation-1
const Sidebar: React.FC = () => {
  const menus = [
    {
      icon: <HomeIcon className="h-5 w-5" />,
      title: "Dashboard",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: false,
    },
    {
      icon: <InboxIcon className="h-5 w-5" />,
      title: "Inbox",
      linkTo: "#",
      isNeedBadge: true,
      badge: (
        <span className="ml-auto rounded-full bg-indigo-50 px-2 py-0.5 text-xs font-medium tracking-wide text-indigo-500">
          New
        </span>
      ),
      isSeparator: false,
    },
    {
      icon: <ChatBubbleBottomCenterTextIcon className="h-5 w-5" />,
      title: "Messages",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: false,
    },
    {
      icon: <BellIcon className="h-5 w-5" />,
      title: "Notifications",
      linkTo: "#",
      isNeedBadge: true,
      badge: (
        <span className="ml-auto rounded-full bg-red-50 px-2 py-0.5 text-xs font-medium tracking-wide text-red-500">
          1.2k
        </span>
      ),
      isSeparator: false,
    },
    {
      icon: null,
      title: "Tasks",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: true,
    },
    {
      icon: <DocumentTextIcon className="h-5 w-5" />,
      title: "Available Tasks",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: false,
    },
    {
      icon: <UsersIcon className="h-5 w-5" />,
      title: "Clients",
      linkTo: "#",
      isNeedBadge: true,
      badge: (
        <span className="ml-auto rounded-full bg-green-50 px-2 py-0.5 text-xs font-medium tracking-wide text-green-500">
          15
        </span>
      ),
      isSeparator: false,
    },
    {
      icon: null,
      title: "Settings",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: true,
    },
    {
      icon: <UsersIcon className="h-5 w-5" />,
      title: "Profile",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: false,
    },
    {
      icon: <Cog8ToothIcon className="h-5 w-5" />,
      title: "Settings",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: false,
    },
    {
      icon: <ArrowRightStartOnRectangleIcon className="h-5 w-5" />,
      title: "Logout",
      linkTo: "#",
      isNeedBadge: false,
      badge: null,
      isSeparator: false,
    },
  ];

  const Separator: React.FC<{ title: string }> = ({ title }) => (
    <li className="px-5">
      <div className="flex h-8 items-center">
        <div className="text-sm font-light tracking-wide text-gray-500">
          {title}
        </div>
      </div>
    </li>
  );

  return (
    <aside className="hidden fixed left-0 top-0 h-full flex-col border-r bg-white lg:flex lg:w-64">
      <div className="flex h-16 items-center justify-center">
        <Link to={PATH_DASHBOARD}>
          <img src="/logo.png" alt="logo" className="h-14" />
        </Link>
      </div>
      <div className="overflow-y-hidden hover:overflow-y-auto [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:bg-transparent hover:[&::-webkit-scrollbar-thumb]:bg-gray-300 [&::-webkit-scrollbar-track]:rounded-full [&::-webkit-scrollbar-track]:bg-gray-100 [&::-webkit-scrollbar]:w-2">
        <ul className="flex flex-col space-y-1">
          {menus.map((menu, index) =>
            menu.isSeparator ? (
              <Separator key={index} title={menu.title} />
            ) : (
              <li key={index}>
                <Link
                  to={menu.linkTo}
                  className="relative flex h-11 items-center border-l-4 border-transparent pr-6 text-gray-600 hover:border-indigo-500 hover:bg-gray-50 hover:text-gray-800 focus:outline-none"
                >
                  <span className="ml-4 inline-flex items-center justify-center">
                    {menu.icon}
                  </span>
                  <span className="ml-2 truncate text-sm tracking-wide">
                    {menu.title}
                  </span>
                  {menu.badge}
                </Link>
              </li>
            ),
          )}
        </ul>
      </div>
    </aside>
  );
};

export default Sidebar;
