import { type DialogRoot } from "@base-ui/react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogHeader,
  DialogFooter,
} from "~/components/ui/dialog";
import React, { use, useCallback, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { useGetEventQuery } from "~/hooks";
import { useListEligibleUsersForEventQuery } from "~/hooks/event-registrations";
import { Input } from "~/components/ui/input";
import { InputField } from "~/components/forms";
import { Checkbox } from "~/components/ui/checkbox";
import { useDebounce } from "~/hooks/use-debounce";
import { Button } from "~/components/ui/button";

interface Props {}

export const EnrollStudentDialog: React.FC<Props> = () => {
  const navigate = useNavigate();

  const { id } = useParams();

  const [searchTerm, setSearchTerm] = useState("");
  const [selectedUser, setSelectedUser] = useState<string | null>(null);
  const debouncedSearch = useDebounce(searchTerm, 100);
  const listRef = useRef(null);

  const { data: event, isPending } = useGetEventQuery(id!);
  const { data: eligibleUsers } = useListEligibleUsersForEventQuery({
    eventId: id!,
    searchTerm: debouncedSearch,
  });

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
      if (!["Down", "Up"].includes(event.key)) return;

      setSelectedUser((userId) => {
        if (!collectionRef.current) return userId;

        const activeIndex = userId ? collectionRef.current.findIndex((u) => u.id === userId) : null;

        const newIndex = (() => {
          switch (event.key) {
            case "Down":
              if (activeIndex === collectionRef.current.length - 1) return null;

              if (activeIndex === null) {
                return 0;
              }

              return activeIndex + 1;

            case "Up":
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

  if (isPending) return null;

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <form>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Enroll student for event &ldquo;{event?.titleEn}&rdquo;</DialogTitle>
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
        </DialogContent>
      </form>
    </Dialog>
  );
};
