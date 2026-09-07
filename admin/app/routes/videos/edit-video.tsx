import React, { useCallback, useEffect, useState } from "react";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { Notification } from "~/components/notification";
import { PageTitle } from "~/components/page-title";
import { Button } from "~/components/ui/button";
import { DetailsTable, DetailsTableField as Field } from "~/components/ui/details-table";
import { useGetVideoQuery, useUpdateVideoMutation } from "~/hooks";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { formatInstant } from "~/lib/time";

import { FormFields } from "./form-fields";
import {
  formatDuration,
  sourceLabel,
  toDateInputValue,
  toUpdateVideoInput,
  type VideoFormValues,
} from "./interfaces";

interface Props {}

export const EditVideo: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: video, isPending, isError } = useGetVideoQuery(id!);
  const form = useForm<VideoFormValues>({
    defaultValues: {
      titleEn: "",
      titlePl: "",
      slug: "",
      descriptionEn: "",
      descriptionPl: "",
      recordedOn: "",
      isPublic: false,
      hostId: "",
      youtubeId: "",
    },
  });
  const mutation = useUpdateVideoMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (isPending || !video) return;

    form.reset({
      titleEn: video.titleEn,
      titlePl: video.titlePl,
      slug: video.slug,
      descriptionEn: video.descriptionEn ?? "",
      descriptionPl: video.descriptionPl ?? "",
      recordedOn: toDateInputValue(video.recordedOn),
      isPublic: video.isPublic,
      hostId: video.hostId ?? "",
      youtubeId: video.youtubeId ?? "",
    });
  }, [video, isPending, form]);

  const onSubmit = useCallback(
    async (values: VideoFormValues) => {
      setFormError(null);
      try {
        await mutation.mutateAsync({ id: id!, params: toUpdateVideoInput(values) });
        navigate("/videos");
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            form.setError(field as Path<VideoFormValues>, { message });
          }
          setFormError("Please fix the errors below.");
          return;
        }
        setFormError("Something went wrong updating the video. Please try again.");
      }
    },
    [mutation, id, navigate, form],
  );

  return (
    <AdminLayout title="Edit video">
      <BackButton href="/videos">Back to list</BackButton>
      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !video ? (
        <p className="text-destructive">Video not found.</p>
      ) : (
        <FormProvider {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-6">
            <PageTitle subtitle="Edit video">{video.titleEn}</PageTitle>
            {formError ? (
              <Notification
                title="An error has prevented this video from being saved."
                variant="destructive"
              >
                {formError}
              </Notification>
            ) : null}

            <FormFields provider={video.provider} />

            <div className="flex gap-2">
              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending ? "Updating…" : "Update video"}
              </Button>
            </div>

            {/* Owned by the import script (aside from the YouTube ID, editable
                above), shown so the record can be matched against its source
                without leaving the page. */}
            <section className="grid max-w-2xl gap-2">
              <h3 className="text-lg font-semibold">Metadata</h3>
              <DetailsTable>
                <Field label="Provider">{video.provider}</Field>
                <Field label="Duration">{formatDuration(video.durationSeconds)}</Field>
                <Field label="Created at">{formatInstant(video.insertedAt)}</Field>
                <Field label="Updated at">{formatInstant(video.updatedAt)}</Field>
              </DetailsTable>
            </section>

            {/* video_sources rows, distinct from the metadata above: this is
                what the player actually streams from. */}
            {video.provider === "cloudfront" && (video.sources ?? []).length > 0 ? (
              <section className="grid max-w-2xl gap-2">
                <h3 className="text-lg font-semibold">Sources</h3>
                <DetailsTable>
                  {(video.sources ?? []).map((source) => (
                    <Field key={source.id} label={sourceLabel(source.contentType)} copy monospace>
                      {source.objectKey}
                    </Field>
                  ))}
                </DetailsTable>
              </section>
            ) : null}
          </form>
        </FormProvider>
      )}
    </AdminLayout>
  );
};
