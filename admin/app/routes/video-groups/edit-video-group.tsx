import React, { useCallback, useEffect, useState } from "react";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { Notification } from "~/components/notification";
import { PageTitle } from "~/components/page-title";
import { Button } from "~/components/ui/button";
import {
  useGetVideoGroupQuery,
  useListVideosInVideoGroupQuery,
  useReplaceVideosInVideoGroupMutation,
  useUpdateVideoGroupMutation,
} from "~/hooks";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { FormFields } from "./form-fields";
import { type VideoGroupFormValues, toPatchVideoGroupInput } from "./interfaces";

interface Props {}

export const EditVideoGroup: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: group, isPending, isError } = useGetVideoGroupQuery(id);
  const { data: videos, isPending: videosPending } = useListVideosInVideoGroupQuery(id);
  const form = useForm<VideoGroupFormValues>({
    defaultValues: { isFree: true, currency: "PLN", videoIds: [] },
  });
  const mutation = useUpdateVideoGroupMutation();
  const replaceVideos = useReplaceVideosInVideoGroupMutation();
  const navigate = useNavigate();
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    if (isPending || videosPending || !group) return;

    form.reset({
      titleEn: group.titleEn,
      titlePl: group.titlePl,
      slug: group.slug,
      isFree: !group.isPremium,
      price: group.price ?? undefined,
      currency: group.currency ?? "PLN",
      videoIds: (videos ?? []).map(({ id }) => id),
    });
  }, [group, videos, isPending, videosPending]);

  const onSubmit = useCallback(
    async (values: VideoGroupFormValues) => {
      setFormError(null);
      try {
        await mutation.mutateAsync({ id: id!, params: toPatchVideoGroupInput(values) });
        // Membership and order live behind their own endpoint, because positions
        // are unique per group and are replaced wholesale.
        await replaceVideos.mutateAsync({ id: id!, videoIds: values.videoIds ?? [] });
        navigate("/videos");
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            form.setError(field as Path<VideoGroupFormValues>, { message });
          }
          setFormError("Please fix the errors below.");
          return;
        }
        setFormError("Something went wrong saving the series. Please try again.");
      }
    },
    [mutation, id],
  );

  return (
    <AdminLayout title="Edit video series">
      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !group ? (
        <p className="text-desctructive">Video series not found.</p>
      ) : (
        <FormProvider {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="grid gap-3">
            <BackButton href="/videos">Back to video series</BackButton>
            <PageTitle subtitle="Edit video series">{group.titleEn}</PageTitle>
            {formError ? (
              <Notification
                title="An error has prevented this series from being saved."
                variant="destructive"
              >
                {formError}
              </Notification>
            ) : null}

            <FormFields showVideos />

            <div className="flex gap-2">
              <Button type="submit" disabled={mutation.isPending || replaceVideos.isPending}>
                {mutation.isPending || replaceVideos.isPending ? "Updating…" : "Update series"}
              </Button>
            </div>
          </form>
        </FormProvider>
      )}
    </AdminLayout>
  );
};
