import { Input } from "~/components/ui/input";
import {
  CalendarDotsIcon as CalendarDots,
  UsersIcon as Users,
  VideoCameraIcon as VideoCamera,
  ArticleIcon,
} from "@phosphor-icons/react";
import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { NavLink, useLocation, useNavigate } from "react-router";

import { NavUser } from "~/components/nav-user";
import { ThemeToggle } from "~/components/theme-toggle";
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
  { title: "Blog", to: "/blog-posts", icon: ArticleIcon },
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
  searchFormAction?: string;
}

export function AdminLayout({ title, children, searchFormAction }: Props) {
  const { pathname } = useLocation();
  const { data: session, isLoading } = useGetSessionQuery();
  const navigate = useNavigate();
  const inputRef = useRef<HTMLInputElement>(null);
  const [searchHasFocus, setSearchHasFocus] = useState(false);

  useEffect(() => {
    if (isLoading || session) return;
    const qs = new URLSearchParams({ ref: location.pathname });
    location.href = `/sign-in?${qs}`;
  }, [isLoading, session]);

  const onSubmit: React.SubmitEventHandler = useCallback(
    (e) => {
      e.preventDefault();
      const searchTerm = inputRef.current?.value;
      const qs = new URLSearchParams({ q: searchTerm ?? "" });
      const nextUrl = `${searchFormAction}?${qs}`;
      navigate(nextUrl);
    },
    [searchFormAction, navigate, inputRef],
  );

  const onKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (!searchFormAction || !inputRef.current) return;

      const hasFocus = document.activeElement === inputRef.current;

      switch (e.key) {
        case "/": {
          if (hasFocus) return;

          e.preventDefault();
          inputRef.current.focus();
          inputRef.current.select();
          return;
        }

        case "Escape": {
          if (!hasFocus) return;

          e.preventDefault();
          inputRef.current.blur();
        }
      }
    },
    [inputRef, searchFormAction],
  );

  const onSearchInputFocus: React.FocusEventHandler<HTMLInputElement> = useCallback(
    (e) => {
      setSearchHasFocus(document.activeElement === e.currentTarget);
    },
    [setSearchHasFocus],
  );

  useEffect(() => {
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [onKeyDown]);

  if (isLoading) {
    return "Loading";
  }

  return (
    <SidebarProvider defaultOpen={sidebarDefaultOpen()} className="overflow-x-hidden">
      <title>{title ? `${title} | Homeo sapiens` : "Homeo sapiens"}</title>
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
                <SidebarMenuItem>
                  <SidebarMenuButton render={<a href="/" />}>
                    <span>Back to website</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>{session && <NavUser user={session} />}</SidebarFooter>
      </Sidebar>
      <SidebarInset className="overflow-x-hidden">
        <header className="flex h-14 items-center gap-2 border-b px-4">
          <SidebarTrigger />
          <Separator orientation="vertical" className="mx-2" />
          {searchFormAction ? (
            <form onSubmit={onSubmit}>
              <Input
                name="q"
                placeholder={searchHasFocus ? "Search..." : "Press / to search..."}
                ref={inputRef}
                onFocus={onSearchInputFocus}
                onBlur={onSearchInputFocus}
              />
              <input type="submit" className="sr-only" />
            </form>
          ) : null}
          <ThemeToggle className="ml-auto" />
        </header>
        <main className="flex-1 p-4">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  );
}
