import { CaretLeftIcon as CaretLeft } from "@phosphor-icons/react";
import { Link, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Button, buttonVariants } from "~/components/ui/button";
import { useGetEventQuery } from "~/hooks";
import { DetailsTable, DataTableField as Field } from "~/components/ui/details-table";
import Markdown from "react-markdown";
import { PencilIcon } from "@phosphor-icons/react/ssr";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { PageTitle } from "~/components/page-title";

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
        <Link
          to="/events"
          className={buttonVariants({ size: "sm", variant: "ghost", className: "w-fit" })}
        >
          <CaretLeft />
          Back to events
        </Link>

        {isPending ? (
          <p className="text-muted-foreground">Loading…</p>
        ) : isError || !event ? (
          <p className="text-destructive">Event not found.</p>
        ) : (
          <>
            <header className="flex justify-between">
              <PageTitle subtitle="Event details">{event.titleEn}</PageTitle>
              <div>
                <Link
                  to={`/events/${event.id}/edit`}
                  className={buttonVariants({ variant: "secondary" })}
                >
                  <PencilIcon className="w-5" />
                  Edit
                </Link>
              </div>
            </header>

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
            <section className="mt-6">
              <h4 className="text-xl font-bold">Descriptions</h4>
              <div className="mt-6 flex gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Polish</CardTitle>
                  </CardHeader>
                  <CardContent className="prose w-[65ch] text-foreground">
                    <Markdown>{event.descriptionPl}</Markdown>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader>
                    <CardTitle>English</CardTitle>
                  </CardHeader>
                  <CardContent className="prose w-[65ch] text-foreground">
                    <Markdown>{event.descriptionEn}</Markdown>
                  </CardContent>
                </Card>
              </div>
            </section>
          </>
        )}
      </div>
    </AdminLayout>
  );
}
