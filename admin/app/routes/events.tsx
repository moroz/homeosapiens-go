import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { useListEventsQuery } from "~/hooks";

const dateStyle = { dateStyle: "medium", timeStyle: "short" } as const;

function formatInstant(iso: string) {
  return Temporal.Instant.from(iso).toLocaleString("en-GB", dateStyle);
}

function formatRange(startsAt: string, endsAt: string) {
  return `${formatInstant(startsAt)} – ${formatInstant(endsAt)}`;
}

export default function Events() {
  const { data, isPending, isError } = useListEventsQuery();
  const events = data?.data ?? [];

  return (
    <div className="flex flex-col gap-4">
      <h1 className="font-heading text-2xl font-semibold">Events</h1>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Title (EN)</TableHead>
            <TableHead>Title (PL)</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Format</TableHead>
            <TableHead>When</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isPending ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center text-muted-foreground">
                Loading…
              </TableCell>
            </TableRow>
          ) : isError ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center text-destructive">
                Failed to load events.
              </TableCell>
            </TableRow>
          ) : events.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center text-muted-foreground">
                No events yet.
              </TableCell>
            </TableRow>
          ) : (
            events.map((event) => (
              <TableRow key={event.id}>
                <TableCell>{event.titleEn}</TableCell>
                <TableCell>{event.titlePl}</TableCell>
                <TableCell>{event.eventType}</TableCell>
                <TableCell>{event.isVirtual ? "Virtual" : "In person"}</TableCell>
                <TableCell>{formatRange(event.startsAt, event.endsAt)}</TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
