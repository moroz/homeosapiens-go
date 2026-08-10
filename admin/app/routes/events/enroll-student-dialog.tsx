import { type DialogRoot } from "@base-ui/react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogHeader,
  DialogFooter,
} from "~/components/ui/dialog";
import React, { useCallback, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { useGetEventQuery } from "~/hooks";
import {
  useEnrollStudentForEventMutation,
  useListEligibleUsersForEventQuery,
} from "~/hooks/event-registrations";
import { InputField } from "~/components/forms";
import { useDebounce } from "~/hooks/use-debounce";
import { Button } from "~/components/ui/button";

interface Props {}

export const EnrollStudentDialog: React.FC<Props> = () => {
  const navigate = useNavigate();

  const { id: eventId } = useParams();

  const [searchTerm, setSearchTerm] = useState("");
  const [selectedUser, setSelectedUser] = useState<string | null>(null);
  const debouncedSearch = useDebounce(searchTerm, 100);
  const listRef = useRef(null);

  const { data: event, isPending } = useGetEventQuery(eventId!);
  const { data: eligibleUsers } = useListEligibleUsersForEventQuery({
    eventId: eventId!,
    searchTerm: debouncedSearch,
  });

  const mutation = useEnrollStudentForEventMutation();

  const collectionRef = useRef<typeof eligibleUsers>([]);
  collectionRef.current = eligibleUsers;

  const onSearchChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (event) => {
      setSearchTerm(event.currentTarget.value);
    },
    [setSearchTerm],
  );

  const close = useCallback(() => navigate(".."), [navigate]);

  const onOpenChange = useCallback(
    (open: boolean, _eventDetails: DialogRoot.ChangeEventDetails) => {
      if (!open) close();
    },
    [navigate],
  );

  const onSelectedUserChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (event) => {
      setSelectedUser(event.currentTarget.value);
    },
    [setSelectedUser],
  );

  const onKeyDown: React.KeyboardEventHandler<HTMLInputElement> = useCallback(
    (event) => {
      if (!["ArrowDown", "ArrowUp"].includes(event.key)) return;

      setSelectedUser((userId) => {
        if (!collectionRef.current) return userId;

        const activeIndex = userId ? collectionRef.current.findIndex((u) => u.id === userId) : null;

        const newIndex = (() => {
          switch (event.key) {
            case "ArrowDown":
              if (activeIndex === collectionRef.current.length - 1) return null;

              if (activeIndex === null) {
                return 0;
              }

              return activeIndex + 1;

            case "ArrowUp":
              if (activeIndex === 0 || activeIndex === null) return null;
              return activeIndex - 1;

            default:
              return activeIndex;
          }
        })();

        if (newIndex === null) return null;
        return collectionRef.current[newIndex]?.id;
      });
    },
    [listRef, collectionRef],
  );

  const onSubmit: React.SubmitEventHandler<HTMLFormElement> = useCallback(
    async (e) => {
      e.preventDefault();
      if (!selectedUser) return;

      await mutation.mutateAsync({ eventId: eventId!, userId: selectedUser });
      close();
    },
    [mutation, selectedUser, eventId, close],
  );

  if (isPending) return null;

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={onSubmit} className="space-y-3">
          <DialogHeader>
            <DialogTitle className="leading-normal">
              Enroll a student for event:
              <br />
              {event?.titleEn}
            </DialogTitle>
          </DialogHeader>
          <div className="flex flex-col gap-2">
            <InputField
              name="q"
              label="Search users:"
              onKeyDown={onKeyDown}
              value={searchTerm}
              autoComplete="off"
              onChange={onSearchChange}
              autoFocus
            />
            <fieldset
              className="flex h-48 flex-col gap-1 overflow-y-auto rounded-md border border-input py-1"
              ref={listRef}
            >
              {eligibleUsers?.map((user) => {
                return (
                  <label
                    key={user.id}
                    className="flex cursor-pointer items-center gap-2 px-2 py-1 has-checked:bg-primary has-checked:text-white has-focus-visible:ring-3 has-focus-visible:ring-ring/50"
                  >
                    <input
                      type="radio"
                      name="userId"
                      value={user.id}
                      className="sr-only"
                      checked={selectedUser === user.id}
                      onChange={onSelectedUserChange}
                    />
                    {user.givenName} {user.familyName} ({user.email})
                    {user.enrolled && <span>(already enrolled)</span>}
                  </label>
                );
              })}
            </fieldset>
          </div>
          <DialogFooter>
            <Button variant="default" type="submit">
              Enroll student
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};
