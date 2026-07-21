import {
  CalendarDotsIcon as CalendarDots,
  UsersIcon as Users,
  VideoCameraIcon as VideoCamera,
} from "@phosphor-icons/react";
import { useEffect, type ReactNode } from "react";
import { Helmet } from "react-helmet-async";
import { NavLink, useLocation } from "react-router";

import { NavUser } from "~/components/nav-user";
import { Separator } from "~/components/ui/separator";
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

/**
 * The SPA has no server render, so `SidebarProvider` never seeds its initial
 * open state from the persisted cookie. Read it here so remounting the layout
 * (it renders per page) keeps the collapsed/expanded state the user chose.
 */
function sidebarDefaultOpen() {
  const match = document.cookie.match(/(?:^|; )sidebar_state=([^;]+)/);
  return match ? match[1] === "true" : true;
}

interface Props {
  title?: ReactNode;
  children?: ReactNode;
}

export function AdminLayout({ title, children }: Props) {
  const { pathname } = useLocation();
  const { data: session, isLoading } = useGetSessionQuery();

  useEffect(() => {
    if (isLoading || session) return;
    const qs = new URLSearchParams({ ref: location.pathname });
    location.href = `/sign-in?${qs}`;
  }, [isLoading, session]);

  if (isLoading) {
    return "Loading";
  }

  return (
    <SidebarProvider defaultOpen={sidebarDefaultOpen()}>
      <title>{`${title} | Homeo sapiens`}</title>
      <Sidebar>
        <SidebarHeader>
          <h1 className="px-2 py-1.5 font-heading text-lg font-semibold">Homeo sapiens</h1>
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
          <Separator orientation="vertical" className="mx-2" />
        </header>
        <main className="flex-1 p-4">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  );
}
