import React from "react";
import { Link, useParams } from "react-router";
import { useGetVideoQuery } from "~/hooks";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { BackButton } from "~/components/back-button";
import { DetailsTable, DetailsTableField as Field } from "~/components/ui/details-table";
import { buttonVariants } from "~/components/ui/button";
import { PencilIcon } from "@phosphor-icons/react";
import Markdown from "react-markdown";
import { formatDate } from "~/lib/time";

interface Props {}

export const VideoDetails: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: video, isPending, isError } = useGetVideoQuery(id!);

  return (
    <AdminLayout title="Video details">
      <BackButton href="/videos">Back to list</BackButton>
      {isPending ? (
        <p>Loading...</p>
      ) : isError || !video ? (
        <p>Not found</p>
      ) : (
        <div className=" grid gap-4">
          <div className="flex items-start justify-between">
            <PageTitle subtitle="Video details">{video.titleEn}</PageTitle>
            <Link to={`/videos/${video.id}/edit`} className={buttonVariants({ variant: "outline" })}>
              <PencilIcon className="w-5" />
              Edit
            </Link>
          </div>
          <DetailsTable>
            <Field label="ID" copy monospace>
              {video.id}
            </Field>
            <Field label="Title (EN)">{video.titleEn}</Field>
            <Field label="Title (PL)">{video.titlePl}</Field>
            <Field label="Slug" monospace copy>
              {video.slug}
            </Field>
            <Field label="Provider">{video.provider}</Field>
            <Field label="Recorded on">{formatDate(video.recordedOn)}</Field>
            <Field label="Description (EN)">
              {video.descriptionEn && <Markdown>{video.descriptionEn}</Markdown>}
            </Field>
            <Field label="Description (PL)">
              {video.descriptionPl && <Markdown>{video.descriptionPl}</Markdown>}
            </Field>
          </DetailsTable>
        </div>
      )}
    </AdminLayout>
  );
};
