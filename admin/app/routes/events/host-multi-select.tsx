import { useState } from "react";
import { Checkbox } from "~/components/ui/checkbox";
import { Input } from "~/components/ui/input";
import { useListHostsQuery } from "~/hooks";

export function HostMultiSelect({
  value,
  onChange,
}: {
  value: string[];
  onChange: (ids: string[]) => void;
}) {
  const [search, setSearch] = useState("");
  const { data: hosts, isPending } = useListHostsQuery();

  const filtered = (hosts?.data ?? []).filter((host) =>
    `${host.givenName} ${host.familyName}`.toLowerCase().includes(search.toLowerCase()),
  );

  function toggle(id: string) {
    onChange(value.includes(id) ? value.filter((existing) => existing !== id) : [...value, id]);
  }

  return (
    <div className="flex flex-col gap-2">
      <Input
        placeholder="Search hosts…"
        value={search}
        onChange={(event) => setSearch(event.target.value)}
      />
      <div className="flex max-h-48 flex-col gap-1 overflow-y-auto rounded-md border border-input p-2">
        {isPending ? (
          <p className="text-sm text-muted-foreground">Loading hosts…</p>
        ) : filtered.length === 0 ? (
          <p className="text-sm text-muted-foreground">No hosts found.</p>
        ) : (
          filtered.map((host) => (
            <label key={host.id} className="flex items-center gap-2 rounded-sm p-1 text-sm">
              <Checkbox checked={value.includes(host.id)} onCheckedChange={() => toggle(host.id)} />
              {host.givenName} {host.familyName}
            </label>
          ))
        )}
      </div>
    </div>
  );
}
