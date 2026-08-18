import React, { useCallback, useMemo } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { Link, useParams } from "react-router";
import {
  useGetBlogPostQuery,
  usePublishBlogPostMutation,
  useUnpublishBlogPostMutation,
} from "~/hooks";
import { DetailsTable, DetailsTableField as Field } from "~/components/ui/details-table";
import { formatInstant } from "~/lib/time";
import { BackButton } from "~/components/back-button";
import Markdown from "react-markdown";
import { Button, buttonVariants } from "~/components/ui/button";
import { EyeSlashIcon, GlobeIcon, PaperPlaneTiltIcon, PencilIcon } from "@phosphor-icons/react";
import { Badge } from "~/components/ui/badge";

interface Props {}

export const BlogPostDetails: React.FC<Props> = ({}: Props) => {
  const { id } = useParams();
  const { data: blogPost, isPending, isError } = useGetBlogPostQuery(id!);
  const publish = usePublishBlogPostMutation();
  const unpublish = useUnpublishBlogPostMutation();

  const canBePublished = useMemo(() => Boolean(blogPost?.body), [blogPost]);

  const onClickPublish = useCallback(async () => {
    if (!canBePublished) return;

    await publish.mutateAsync(id!);
  }, [publish, id, canBePublished]);

  const onClickUnpublish = useCallback(async () => {
    if (!confirm("Are you sure you want to unpublish this blog post?")) return;

    await unpublish.mutateAsync(id!);
  }, [unpublish, id]);

  return (
    <AdminLayout title="Blog post details">
      <BackButton href="/blog-posts">Back to list</BackButton>
      {isPending ? (
        <p>Loading...</p>
      ) : isError || !blogPost ? (
        <p>No post found.</p>
      ) : (
        <div className="grid gap-4">
          <div className="flex items-start justify-between">
            <PageTitle subtitle="Blog post details">{blogPost.title}</PageTitle>
            <div className="grid justify-end gap-1">
              <div className="flex items-center justify-end gap-3">
                {blogPost.publishedAt ? (
                  <>
                    <a
                      href={`/blog/${blogPost.slug}`}
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
                  to={`/blog-posts/${blogPost.id}/edit`}
                  className={buttonVariants({ variant: "outline" })}
                >
                  <PencilIcon className="w-5" />
                  Edit
                </Link>
              </div>
              {!canBePublished && (
                <p className="text-sm text-muted-foreground">
                  Please add a body before publishing this post.
                </p>
              )}
            </div>
          </div>
          <DetailsTable>
            <Field label="ID" copy>
              {blogPost.id}
            </Field>
            <Field label="Title">{blogPost.title}</Field>
            <Field label="Slug" className="font-mono" copy>
              {blogPost.slug}
            </Field>
            <Field label="Published at">
              {blogPost.publishedAt ? (
                formatInstant(blogPost.publishedAt)
              ) : (
                <Badge variant="secondary">draft</Badge>
              )}
            </Field>
            <Field label="Created at">{formatInstant(blogPost.insertedAt)}</Field>
            <Field label="Body">
              {blogPost.body ? (
                <div className="prose p-2 dark:prose-invert">
                  <Markdown>{blogPost.body}</Markdown>
                </div>
              ) : (
                <span className="text-muted-foreground">(No data)</span>
              )}
            </Field>
          </DetailsTable>
        </div>
      )}
    </AdminLayout>
  );
};
