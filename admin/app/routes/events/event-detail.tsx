import { Link, Outlet, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Button, buttonVariants } from "~/components/ui/button";
import { useGetEventQuery, usePublishEventMutation, useUnpublishEventMutation } from "~/hooks";
import { DetailsTableField as Field, DetailsTable } from "~/components/ui/details-table";
import { GlobeIcon, PencilIcon, UserListIcon } from "@phosphor-icons/react/ssr";
import { PageTitle } from "~/components/page-title";
import { formatInstant } from "~/lib/time";
import { BackButton } from "~/components/back-button";
import { Badge } from "~/components/ui/badge";
import { EventDescriptionCard } from "~/routes/events/event-description-card";
import { useCallback, useMemo } from "react";
import { EyeSlashIcon, PaperPlaneTiltIcon } from "@phosphor-icons/react";
import { formatPrice } from "~/lib/money";

export default function EventDetail() {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);
  const publish = usePublishEventMutation();
  const unpublish = useUnpublishEventMutation();

  const onClickUnpublish = useCallback(async () => {
    if (!confirm("Are you sure you want to unpublish this event?")) return;

    await unpublish.mutateAsync(id!);
  }, [unpublish, id]);

  const canBePublished = useMemo(() => {
    return Boolean(event?.descriptionEn && event.descriptionPl);
  }, [event]);

  const onClickPublish = useCallback(async () => {
    if (!canBePublished) return;

    await publish.mutateAsync(id!);
  }, [publish, id, canBePublished]);

  return (
    <AdminLayout title={event?.titleEn ?? "Event"}>
      <div className="flex flex-col gap-3">
        <BackButton href="/events">Back to events</BackButton>

        {isPending ? (
          <p className="text-muted-foreground">Loading…</p>
        ) : isError || !event ? (
          <p className="text-destructive">Event not found.</p>
        ) : (
          <>
            <header className="flex justify-between">
              <PageTitle subtitle="Event details">{event.titleEn}</PageTitle>
              <div className="grid justify-end gap-1">
                <div className="flex items-center justify-end gap-3">
                  <Link
                    to={`/events/${event.id}/attendants`}
                    className={buttonVariants({ variant: "outline" })}
                  >
                    <UserListIcon className="w-5" />
                    Enrolled students
                  </Link>
                  {event.publishedAt ? (
                    <>
                      <a
                        href={`/events/${event.slug}`}
                        target="_blank"
                        className={buttonVariants({ variant: "outline" })}
                      >
                        <GlobeIcon className="w-5" />
                        View on website
                      </a>
                      <Button variant="destructive" onClick={onClickUnpublish} type="button">
                        <EyeSlashIcon className="w-5" />
                        Unpublish
                      </Button>
                    </>
                  ) : (
                    <Button onClick={onClickPublish} type="button" disabled={!canBePublished}>
                      <PaperPlaneTiltIcon className="w-5" />
                      Publish
                    </Button>
                  )}
                  <Link
                    to={`/events/${event.id}/edit`}
                    className={buttonVariants({ variant: "outline" })}
                  >
                    <PencilIcon className="w-5" />
                    Edit
                  </Link>
                </div>
                {!canBePublished && (
                  <p className="text-sm text-muted-foreground">
                    Please fill in the event description in all languages before publishing this
                    event.
                  </p>
                )}
              </div>
            </header>

            <DetailsTable>
              <Field label="ID" className="font-mono select-all" copy>
                {event.id}
              </Field>
              <Field label="Published at">
                {event.publishedAt ? (
                  formatInstant(event.publishedAt)
                ) : (
                  <div className="grid py-1 text-sm text-muted-foreground">
                    <Badge variant="secondary" className="mb-1">
                      draft
                    </Badge>
                    <p>
                      This event is a draft. It is not visible on the public website. Students
                      cannot sign up for this event.
                    </p>
                    {canBePublished ? (
                      <p>This event has descriptions in all languages and can be published.</p>
                    ) : (
                      <p>
                        You must provide a description in all languages before publishing this
                        event.
                      </p>
                    )}
                  </div>
                )}
              </Field>
              <Field label="Title (PL)">{event.titlePl}</Field>
              <Field label="Subtitle (PL)">{event.subtitlePl}</Field>
              <Field label="Title (EN)">{event.titleEn}</Field>
              <Field label="Subtitle (EN)">{event.subtitleEn}</Field>
              <Field label="Price">
                {event.isFree ? "Free" : formatPrice(event.price!, event.currency!)}
              </Field>
              <Field label="When">
                {formatInstant(event.startsAt)}&ndash;{formatInstant(event.endsAt)}
              </Field>
              <Field label="Slug" className="font-mono select-all" copy>
                {event.slug}
              </Field>
              <Field label="Zoom link" copy className="font-mono select-all">
                {event.meetingUrl}
              </Field>
              <Field label="Hosts">
                {event.hosts?.length ? (
                  <ul className="list-disc pl-6">
                    {event.hosts.map((host) => (
                      <li key={host.id}>
                        {host.givenName} {host.familyName}
                      </li>
                    ))}
                  </ul>
                ) : null}
              </Field>
            </DetailsTable>
            <section className="mt-6">
              <h4 className="text-xl font-bold">Descriptions</h4>
              <div className="mt-6 flex max-w-full flex-wrap gap-6">
                <EventDescriptionCard language="pl" event={event} />
                <EventDescriptionCard language="en" event={event} />
              </div>
            </section>
            <Outlet />
          </>
        )}
      </div>
    </AdminLayout>
  );
}
