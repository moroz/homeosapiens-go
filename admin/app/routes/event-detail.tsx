import { CaretLeftIcon as CaretLeft } from "@phosphor-icons/react";
import { Link, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Badge } from "~/components/ui/badge";
import { Button } from "~/components/ui/button";
import { useGetEventQuery } from "~/hooks";

const dateStyle = { dateStyle: "full", timeStyle: "short" } as const;

function formatInstant(iso: string) {
  return Temporal.Instant.from(iso).toLocaleString("en-GB", dateStyle);
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col gap-1">
      <dt className="text-xs font-medium text-muted-foreground">{label}</dt>
      <dd className="text-sm">{children}</dd>
    </div>
  );
}

export default function EventDetail() {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);

  return (
    <AdminLayout title={event?.titleEn ?? "Event"}>
      <div className="flex flex-col gap-6">
        <Button variant="ghost" size="sm" className="w-fit" render={<Link to="/events" />}>
          <CaretLeft />
          Back to events
        </Button>

        {isPending ? (
          <p className="text-muted-foreground">Loading…</p>
        ) : isError || !event ? (
          <p className="text-destructive">Event not found.</p>
        ) : (
          <div className="flex flex-col gap-6">
            <p className="text-muted-foreground">{event.titlePl}</p>

            <dl className="grid grid-cols-2 gap-4 sm:grid-cols-3">
              <Field label="Type">
                <Badge variant="outline">{event.eventType}</Badge>
              </Field>
              <Field label="Format">
                <Badge variant="outline">{event.isVirtual ? "Virtual" : "In person"}</Badge>
              </Field>
              <Field label="Slug">{event.slug}</Field>
              <Field label="Starts">{formatInstant(event.startsAt)}</Field>
              <Field label="Ends">{formatInstant(event.endsAt)}</Field>
            </dl>

            {(event.subtitleEn || event.subtitlePl) && (
              <div className="flex flex-col gap-1">
                {event.subtitleEn && <p className="text-sm">{event.subtitleEn}</p>}
                {event.subtitlePl && (
                  <p className="text-sm text-muted-foreground">{event.subtitlePl}</p>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </AdminLayout>
  );
}
