// place files you want to import through the `$lib` alias in this folder.
import { SquareUser, Home, Settings } from 'lucide-svelte';

export const menus = [
    {
        title: "Dashboard",
        icon: Home,
        link: "/dashboard",
    },
    {
        title: "Profile",
        icon: SquareUser,
        link: "/profile",
    },
    {
        title: "Setting",
        icon: Settings,
        link: "/setting",
    },
];