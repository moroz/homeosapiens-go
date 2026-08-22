import { ArrowDownIcon, ArrowUpIcon, PlusIcon, XIcon } from "@phosphor-icons/react";
import { useMemo, useState } from "react";

import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { useListVideosQuery, type Video } from "~/hooks";

interface Props {
  /** Video IDs in playback order. */
  value: string[];
  onChange: (ids: string[]) => void;
}

/**
 * VideoPicker edits both the membership and the order of a series in one list:
 * positions are unique per group, so the API replaces them wholesale and the UI
 * follows suit rather than pretending each video can be moved on its own.
 */
export function VideoPicker({ value, onChange }: Props) {
  const [search, setSearch] = useState("");
  const { data, isPending } = useListVideosQuery();

  const videosById = useMemo(() => {
    const map = new Map<string, Video>();
    for (const video of data?.data ?? []) map.set(video.id, video);
    return map;
  }, [data]);

  const selected = value.map((id) => videosById.get(id)).filter((v) => v != null);

  const available = (data?.data ?? []).filter(
    (video) =>
      !value.includes(video.id) &&
      `${video.titleEn} ${video.titlePl} ${video.slug}`
        .toLowerCase()
        .includes(search.toLowerCase()),
  );

  function move(index: number, delta: number) {
    const target = index + delta;
    if (target < 0 || target >= value.length) return;

    const next = [...value];
    [next[index], next[target]] = [next[target], next[index]];
    onChange(next);
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-2">
        <p className="text-sm font-medium">In this playlist ({value.length})</p>
        {selected.length === 0 ? (
          <p className="rounded-md border border-input p-3 text-sm text-muted-foreground">
            No videos yet. Add them from the list below.
          </p>
        ) : (
          <ol className="flex flex-col gap-1 rounded-md border border-input p-2">
            {selected.map((video, index) => (
              <li key={video.id} className="flex items-center gap-2 rounded-sm p-1 text-sm">
                <span className="w-6 text-muted-foreground tabular-nums">{index + 1}.</span>
                <span className="flex-1">{video.titleEn}</span>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  aria-label="Move up"
                  disabled={index === 0}
                  onClick={() => move(index, -1)}
                >
                  <ArrowUpIcon />
                </Button>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  aria-label="Move down"
                  disabled={index === selected.length - 1}
                  onClick={() => move(index, 1)}
                >
                  <ArrowDownIcon />
                </Button>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  aria-label="Remove from series"
                  onClick={() => onChange(value.filter((id) => id !== video.id))}
                >
                  <XIcon />
                </Button>
              </li>
            ))}
          </ol>
        )}
      </div>

      <div className="flex flex-col gap-2">
        <Input
          placeholder="Search videos…"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
        />
        <div className="flex max-h-64 flex-col gap-1 overflow-y-auto rounded-md border border-input p-2">
          {isPending ? (
            <p className="text-sm text-muted-foreground">Loading videos…</p>
          ) : available.length === 0 ? (
            <p className="text-sm text-muted-foreground">No videos found.</p>
          ) : (
            available.map((video) => (
              <div key={video.id} className="flex items-center gap-2 rounded-sm p-1 text-sm">
                <span className="flex-1">
                  {video.titleEn}
                  <span className="ml-2 text-muted-foreground">{video.slug}</span>
                </span>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  aria-label="Add to series"
                  onClick={() => onChange([...value, video.id])}
                >
                  <PlusIcon />
                </Button>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
