import React, { useCallback, useState } from "react";
import { useNavigate, useParams } from "react-router";
import Markdown from "react-markdown";

import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import { Textarea } from "~/components/ui/textarea";
import { useGetEventQuery, useUpdateEventMutation, type EventDetails } from "~/hooks";
import type { DialogRoot } from "@base-ui/react";

/** `:locale` route param to the event field it edits. */
const FIELD_BY_LOCALE = { en: "descriptionEn", pl: "descriptionPl" } as const;
const LABEL_BY_LOCALE = { en: "English", pl: "Polish" } as const;

type Locale = keyof typeof FIELD_BY_LOCALE;

const DIRTY_MESSAGE = "You have unsaved changes. Are you sure you want to discard them?";

const isMac = navigator.platform.startsWith("Mac") || navigator.platform === "iPhone";

function isLocale(value: string | undefined): value is Locale {
  return value != null && value in FIELD_BY_LOCALE;
}

/**
 * Modal route (`events/:id/description/:locale`) editing one description of an
 * event in Markdown, opened from the detail view. Saving PATCHes just that
 * field; closing discards the draft.
 */
export const MarkdownEditorDialog: React.FC = () => {
  const { id, locale } = useParams();
  const { data: event, isPending } = useGetEventQuery(id);

  if (!isLocale(locale) || !id) return null;

  return (
    <Editor
      key={`${id}:${locale}`}
      id={id}
      locale={locale}
      initialValue={event?.[FIELD_BY_LOCALE[locale]] ?? ""}
      event={event ?? null}
      isPending={isPending}
    />
  );
};

interface EditorProps {
  id: string;
  locale: Locale;
  initialValue: string;
  isPending: boolean;
  event: EventDetails | null;
}

/**
 * Split out so `initialValue` seeds `useState` once the event has loaded — the
 * `key` on the parent remounts this when the loaded value arrives.
 */
const Editor: React.FC<EditorProps> = ({ id, locale, initialValue, event, isPending }) => {
  const navigate = useNavigate();
  const mutation = useUpdateEventMutation();
  const [draft, setDraft] = useState(initialValue);
  const [dirty, setDirty] = useState(false);

  /** `..` resolves against the route hierarchy, i.e. back to the detail view. */
  const close = () => navigate("..");

  const onChange: React.ChangeEventHandler<HTMLTextAreaElement> = useCallback(
    (e) => {
      setDirty(true);
      setDraft(e.currentTarget.value);
    },
    [setDraft, setDirty],
  );

  const onOpenChange = useCallback(
    (open: boolean, eventDetails: DialogRoot.ChangeEventDetails) => {
      if (eventDetails.reason === "escape-key" && dirty) {
        eventDetails.cancel();
        return;
      }

      if (dirty && !open && !confirm(DIRTY_MESSAGE)) {
        return;
      }

      if (!open) close();
    },
    [dirty, close, open],
  );

  const onKeyDown: React.EventHandler<React.KeyboardEvent> = useCallback(
    async (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "s") {
        e.preventDefault();
        await save();
      }
    },
    [save],
  );

  async function save() {
    try {
      await mutation.mutateAsync({
        id,
        params: { [FIELD_BY_LOCALE[locale]]: draft.trim() === "" ? null : draft },
      });
      setDirty(false);
    } catch {
      // Surfaced below via `mutation.isError`.
    }
  }

  function saveAndClose() {
    save().then(() => close());
  }

  if (!event) return null;

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-7xl" onKeyDown={onKeyDown}>
        <DialogHeader>
          <DialogTitle>
            Edit description for event &ldquo;{event.titleEn}&rdquo; &mdash;{" "}
            {LABEL_BY_LOCALE[locale]}
          </DialogTitle>
          <DialogDescription>
            You can use{" "}
            <a href="https://www.markdownguide.org/" target="_blank" rel="noopener noreferrer">
              Markdown
            </a>
            . The preview approximates the published page.
          </DialogDescription>
        </DialogHeader>
        {isPending ? (
          <p className="text-muted-foreground">Loading&hellip;</p>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2">
            <Textarea
              aria-label="Markdown source"
              className="h-120 font-mono text-sm"
              value={draft}
              onChange={onChange}
              autoFocus
            />
            <div className="prose max-h-120 max-w-none overflow-y-auto rounded-md border p-4 text-foreground dark:prose-invert">
              <Markdown>{draft}</Markdown>
            </div>
          </div>
        )}
        <DialogFooter>
          {mutation.isError && (
            <p className="mr-auto self-center text-destructive">Could not save. Try again.</p>
          )}
          <Button type="button" variant="outline" onClick={close}>
            Close{!dirty && " (Esc)"}
          </Button>
          <Button
            type="button"
            onClick={saveAndClose}
            disabled={!dirty || isPending || mutation.isPending}
          >
            Save ({isMac ? "⌘" : "Ctrl"}-S)
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
