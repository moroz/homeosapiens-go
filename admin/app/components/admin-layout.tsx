import {
  CalendarDotsIcon as CalendarDots,
  UsersIcon as Users,
  VideoCameraIcon as VideoCamera,
} from "@phosphor-icons/react";
import { NavLink, Outlet, useLocation } from "react-router";

import { NavUser } from "~/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "~/components/ui/sidebar";
import { useGetSessionQuery } from "~/hooks";

const NAV_ITEMS = [
  { title: "Videos", to: "/videos", icon: VideoCamera },
  { title: "Users", to: "/users", icon: Users },
  { title: "Events", to: "/events", icon: CalendarDots },
];

export function AdminLayout() {
  const { pathname } = useLocation();
  const { data: session } = useGetSessionQuery();

  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader>
          <div className="px-2 py-1.5 font-heading text-lg font-semibold">Homeo sapiens</div>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Admin</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {NAV_ITEMS.map((item) => (
                  <SidebarMenuItem key={item.to}>
                    <SidebarMenuButton
                      isActive={pathname === item.to || pathname.startsWith(item.to + "/")}
                      tooltip={item.title}
                      render={<NavLink to={item.to} />}
                    >
                      <item.icon />
                      <span>{item.title}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>{session && <NavUser user={session} />}</SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <header className="flex h-14 items-center gap-2 border-b px-4">
          <SidebarTrigger />
        </header>
        <main className="flex-1 p-4">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
