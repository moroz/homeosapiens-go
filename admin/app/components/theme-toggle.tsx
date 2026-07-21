import { CheckIcon, DesktopIcon, MoonIcon, SunIcon, type Icon } from "@phosphor-icons/react";

import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { useTheme, type Theme } from "~/hooks";

const THEME_OPTIONS: { value: Theme; label: string; icon: Icon }[] = [
  { value: "light", label: "Light", icon: SunIcon },
  { value: "dark", label: "Dark", icon: MoonIcon },
  { value: "system", label: "System", icon: DesktopIcon },
];

interface Props {
  className?: string;
}

export function ThemeToggle({ className }: Props) {
  const { theme, setTheme } = useTheme();
  const CurrentIcon = THEME_OPTIONS.find((option) => option.value === theme)?.icon ?? DesktopIcon;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger render={<Button variant="ghost" size="icon" className={className} />}>
        <CurrentIcon />
        <span className="sr-only">Toggle theme</span>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {THEME_OPTIONS.map((option) => (
          <DropdownMenuItem key={option.value} onClick={() => setTheme(option.value)}>
            <option.icon />
            {option.label}
            {option.value === theme && <CheckIcon className="ml-auto" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
