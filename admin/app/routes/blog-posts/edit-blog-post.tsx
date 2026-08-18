import React, { useCallback, useEffect, useState } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { useNavigate, useParams } from "react-router";
import { useGetBlogPostQuery, useUpdateBlogPostMutation } from "~/hooks";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { type BlogPostFormValues, toUpdateBlogPostInput } from "./interfaces";
import { FormFields } from "./form-fields";
import { PageTitle } from "~/components/page-title";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { Notification } from "~/components/notification";
import { Button } from "~/components/ui/button";
import { BackButton } from "~/components/back-button";

interface Props {}

export const EditBlogPost: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: blogPost, isPending, isError } = useGetBlogPostQuery(id!);
  const form = useForm<BlogPostFormValues>({
    defaultValues: { language: "en" },
  });
  const mutation = useUpdateBlogPostMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (isPending || !blogPost) return;

    form.reset({
      title: blogPost.title,
      slug: blogPost.slug,
      language: blogPost.language,
      body: blogPost.body ?? "",
    });
  }, [blogPost, isPending, form]);

  const onSubmit = useCallback(
    async (values: BlogPostFormValues) => {
      setFormError(null);
      try {
        const post = await mutation.mutateAsync({
          id: id!,
          params: toUpdateBlogPostInput(values),
        });
        navigate(`/blog-posts/${post.id}`);
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            form.setError(field as Path<BlogPostFormValues>, { message });
          }
          setFormError("Please fix the errors below.");
          return;
        }
        setFormError("Something went wrong updating the blog post. Please try again.");
      }
    },
    [mutation, id, navigate, form],
  );

  return (
    <AdminLayout title="Edit blog post">
      <BackButton href="/blog-posts">Back to list</BackButton>
      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !blogPost ? (
        <p className="text-destructive">Blog post not found.</p>
      ) : (
        <FormProvider {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-4">
            <PageTitle subtitle="Edit blog post">{blogPost.title}</PageTitle>
            {formError ? (
              <Notification
                title="An error has prevented this blog post from being saved."
                variant="destructive"
              >
                {formError}
              </Notification>
            ) : null}

            <FormFields />

            <div className="flex gap-2">
              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending ? "Updating…" : "Update blog post"}
              </Button>
            </div>
          </form>
        </FormProvider>
      )}
    </AdminLayout>
  );
};
