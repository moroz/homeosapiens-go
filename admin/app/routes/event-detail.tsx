import { CaretLeftIcon as CaretLeft } from "@phosphor-icons/react";
import { Link, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Button } from "~/components/ui/button";
import { useGetEventQuery } from "~/hooks";
import { DetailsTable, DataTableField as Field } from "~/components/ui/details-table";

const dateStyle = { dateStyle: "full", timeStyle: "short" } as const;

function formatInstant(iso: string) {
  return Temporal.Instant.from(iso).toLocaleString("en-GB", dateStyle);
}

export default function EventDetail() {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);

  return (
    <AdminLayout title={event?.titleEn ?? "Event"}>
      <div className="flex flex-col gap-3">
        <Button variant="ghost" size="sm" className="w-fit" render={<Link to="/events" />}>
          <CaretLeft />
          Back to events
        </Button>

        {isPending ? (
          <p className="text-muted-foreground">Loading…</p>
        ) : isError || !event ? (
          <p className="text-destructive">Event not found.</p>
        ) : (
          <>
            <div className="flex flex-col">
              <h2 className="text-2xl font-bold">{event.titleEn}</h2>
              <p className="subtitle text-xl text-muted-foreground">Event details</p>
            </div>

            <DetailsTable>
              <Field label="ID" className="font-mono select-all">
                {event.id}
              </Field>
              <Field label="Title (PL)">{event.titlePl}</Field>
              <Field label="Subtitle (PL)">{event.subtitlePl}</Field>
              <Field label="Title (EN)">{event.titleEn}</Field>
              <Field label="Subtitle (EN)">{event.subtitleEn}</Field>
              <Field label="Slug" className="font-mono select-all">
                {event.slug}
              </Field>
            </DetailsTable>
          </>
        )}
      </div>
    </AdminLayout>
  );
}
